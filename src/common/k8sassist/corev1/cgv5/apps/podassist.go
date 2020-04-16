package apps

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/cmd/exec"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type pods struct {
	k8sClient kubernetes.Interface
	cfg       *types.Config
	namespace string
	pod       v1.PodInterface

	genericclioptions.IOStreams
}

func (p *pods) Create(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.Create(k8sPod)
	if err != nil {
		logs.Error("Create pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) Update(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.Update(k8sPod)
	if err != nil {
		logs.Error("Update pod of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) UpdateStatus(pod *model.Pod) (*model.Pod, error) {
	k8sPod := types.ToK8sPod(pod)
	k8sPod, err := p.pod.UpdateStatus(k8sPod)
	if err != nil {
		logs.Error("Create pod status of %s/%s failed. Err:%+v", pod.Name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(k8sPod)
	return modelPod, nil
}

func (p *pods) Delete(name string) error {
	err := p.pod.Delete(name, nil)
	if err != nil {
		logs.Error("delete pod of %s/%s failed. Err:%+v", name, p.namespace, err)
	}
	return err
}

func (p *pods) Get(name string) (*model.Pod, error) {
	pod, err := p.pod.Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("get pod of %s/%s failed. Err:%+v", name, p.namespace, err)
		return nil, err
	}

	modelPod := types.FromK8sPod(pod)
	return modelPod, nil
}

func (p *pods) List(opts model.ListOptions) (*model.PodList, error) {
	podList, err := p.pod.List(types.ToK8sListOptions(opts))
	if err != nil {
		logs.Error("list pods failed. Err:%+v", err)
		return nil, err
	}

	modelPodList := types.FromK8sPodList(podList)
	return modelPodList, nil
}

func (p *pods) GetLogs(name string, opts *model.PodLogOptions) (io.ReadCloser, error) {
	request := p.pod.GetLogs(name, types.ToK8sPodLogOptions(opts))
	if request == nil {
		err := fmt.Errorf("get pod of %s/%s logs failed, request client is null", name, p.namespace)
		logs.Error("%+v", err)
		return nil, err
	}
	return request.Stream()
}

func (p *pods) Exec(podName, containerName string, cmd []string, ptyHandler model.PtyHandler) error {
	req := p.k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(p.namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.cfg, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: types.ToK8sTerminalSizeQueue(ptyHandler),
		Tty:               true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *pods) Cp(podName, container, src, dest string, cmd []string) error {
	reader, outStream := io.Pipe()
	options := &exec.ExecOptions{
		StreamOptions: exec.StreamOptions{
			IOStreams: genericclioptions.IOStreams{
				In:     nil,
				Out:    outStream,
				ErrOut: p.Out,
			},

			Namespace: p.namespace,
			PodName:   podName,
		},

		// TODO: Improve error messages by first testing if 'tar' is present in the container?
		Command:  cmd,
		Executor: &exec.DefaultRemoteExecutor{},
	}

	go func() {
		defer outStream.Close()
		err := p.execute(options, container)
		cmdutil.CheckErr(err)
	}()
	prefix := getPrefix(src)
	prefix = path.Clean(prefix)
	// remove extraneous path shortcuts - these could occur if a path contained extra "../"
	// and attempted to navigate beyond "/" in a remote filesystem
	prefix = stripPathShortcuts(prefix)
	return p.untarAll(reader, dest, prefix)
}

func (p *pods) execute(options *exec.ExecOptions, container string) error {
	if len(options.Namespace) == 0 {
		options.Namespace = p.namespace
	}

	if len(container) > 0 {
		options.ContainerName = container
	}

	options.Config = p.cfg
	options.PodClient = p.k8sClient.CoreV1()

	// if err := options.Validate(); err != nil {
	// 	return err
	// }

	if err := options.Run(); err != nil {
		return err
	}
	return nil
}

func getPrefix(file string) string {
	// tar strips the leading '/' if it's there, so we will too
	return strings.TrimLeft(file, "/")
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(s string) string {
	newPath := path.Clean(s)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}

func (p *pods) untarAll(reader io.Reader, destDir, prefix string) error {
	// TODO: use compression here?
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		// All the files will start with the prefix, which is the directory where
		// they were located on the pod, we need to strip down that prefix, but
		// if the prefix is missing it means the tar was tempered with.
		// For the case where prefix is empty we need to ensure that the path
		// is not absolute, which also indicates the tar file was tempered with.
		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		// basic file information
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		outFile, err := os.Create(destFileName)
		if err != nil {
			return err
		}
		defer outFile.Close()
		if _, err := io.Copy(outFile, tarReader); err != nil {
			return err
		}
		if err := outFile.Close(); err != nil {
			return err
		}
	}

	return nil
}

func NewPods(k8sClient kubernetes.Interface, cfg *types.Config, namespace string, pod v1.PodInterface) *pods {
	return &pods{
		k8sClient: k8sClient,
		cfg:       cfg,
		namespace: namespace,
		pod:       pod,
	}
}
