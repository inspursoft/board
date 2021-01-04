package apps

import (
	"archive/tar"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/inspursoft/board/src/common/k8sassist/corev1/cgv5/types"
	"github.com/inspursoft/board/src/common/model"

	"github.com/astaxie/beego/logs"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type pods struct {
	k8sClient kubernetes.Interface
	cfg       *types.Config
	namespace string
	pod       v1.PodInterface
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

func (p *pods) ShellExec(podName, containerName string, cmd []string, ptyHandler model.PtyHandler) error {
	podExecOptions := corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}
	streamOptions := remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: types.ToK8sTerminalSizeQueue(ptyHandler),
		Tty:               true,
	}
	return p.generatorExec(podExecOptions, streamOptions, podName, "POST")
}

func (p *pods) CopyFromPodExec(podName, containerName string, cmd []string, outStream *io.PipeWriter) error {
	podExecOptions := corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}
	streamOptions := remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: outStream,
		Stderr: os.Stderr,
		Tty:    false,
	}
	return p.generatorExec(podExecOptions, streamOptions, podName, "GET")
}

func (p *pods) CopyToPodExec(podName, containerName string, cmd []string, reader *io.PipeReader) error {
	podExecOptions := corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     reader != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}
	streamOptions := remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  reader,
		Tty:    false,
	}
	return p.generatorExec(podExecOptions, streamOptions, podName, "POST")
}

func (p *pods) generatorExec(podExecOptions corev1.PodExecOptions, streamOptions remotecommand.StreamOptions,
	podName, SPDYMethod string) error {

	var req *rest.Request
	client := p.k8sClient.CoreV1().RESTClient()

	if strings.Compare(SPDYMethod, "GET") == 0 {
		req = client.Get()
	} else if strings.Compare(SPDYMethod, "POST") == 0 {
		req = client.Post()
	} else {
		return errors.New("SPDYMethod should be 'GET' or 'POST'")
	}

	req = req.Resource("pods").
		Name(podName).
		Namespace(p.namespace).
		SubResource("exec").
		VersionedParams(&podExecOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.cfg, SPDYMethod, req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(streamOptions)
	if err != nil {
		return err
	}

	return nil
}

func (p *pods) CopyFromPod(podName, containerName, src, dest string, cmd []string) error {
	reader, outStream := io.Pipe()

	go func() {
		defer outStream.Close()
		p.CopyFromPodExec(podName, containerName, cmd, outStream)
	}()

	prefix := getPrefix(src)
	prefix = path.Clean(prefix)
	// remove extraneous path shortcuts - these could occur if a path contained extra "../"
	// and attempted to navigate beyond "/" in a remote filesystem
	prefix = stripPathShortcuts(prefix)
	return p.untarAll(reader, dest, prefix)
}

func (p *pods) CopyToPod(podName, containerName, src, dest string) error {
	reader, writer := io.Pipe()

	if dest != "/" && strings.HasSuffix(string(dest[len(dest)-1]), "/") {
		dest = dest[:len(dest)-1]
	}

	err := checkDestinationIsDir(p, podName, containerName, dest)
	if err == nil {
		logs.Info("Destination is directory.")
		dest = dest + "/" + path.Base(src)
	}

	go func() {
		defer writer.Close()
		err = makeTar(src, dest, writer)
	}()

	cmd := []string{"tar", "xf", "-"}
	destDir := path.Dir(dest)
	if len(destDir) > 0 {
		cmd = append(cmd, "-C", destDir)
	}
	err = p.CopyToPodExec(podName, containerName, cmd, reader)

	logs.Info("Copying the content of '%s' directory to '%s/%s/%s:%s' finished", src, p.namespace, podName, containerName, destDir)
	return err
}

// execCommand executes the given command inside the specified container remotely
func (p *pods) execCommand(podName, containerName string, stdinReader io.Reader, cmd []string) (string, error) {
	req := p.k8sClient.CoreV1().RESTClient().Get().
		Resource("pods").
		Name(podName).
		Namespace(p.namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   cmd,
			Stdin:     stdinReader != nil,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(p.cfg, "POST", req.URL())

	if err != nil {
		logs.Error("Creating remote command executor failed: %v", err)
		return "", err
	}

	stdOut := bytes.Buffer{}
	stdErr := bytes.Buffer{}

	logs.Debug("Executing command '%v' in namespace='%s', pod='%s', container='%s'", cmd, p.namespace, podName, containerName)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: bufio.NewWriter(&stdOut),
		Stderr: bufio.NewWriter(&stdErr),
		Stdin:  stdinReader,
		Tty:    false,
	})

	logs.Debug("Command stderr: %s", stdErr.String())
	logs.Debug("Command stdout: %s", stdOut.String())

	if err != nil {
		logs.Info("Executing command failed with: %v", err)
		return "", err
	}

	logs.Debug("Command succeeded.")
	if stdErr.Len() > 0 {
		return "", fmt.Errorf("stderr: %v", stdErr.String())
	}

	return stdOut.String(), nil
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

// isDestRelative returns true if dest is pointing outside the base directory,
// false otherwise.
func isDestRelative(base, dest string) bool {
	relative, err := filepath.Rel(base, dest)
	if err != nil {
		return false
	}
	return relative == "." || relative == stripPathShortcuts(relative)
}

func (p *pods) untarAll(reader io.Reader, destDir, prefix string) error {
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
		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		if !isDestRelative(destDir, destFileName) {
			logs.Info("warning: file %s is outside target destination, skipping\n", destFileName)
			continue
		}

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

		if mode&os.ModeSymlink != 0 {
			logs.Info("warning: skipping symlink: %s -> %s\n", destFileName, header.Linkname)
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

// checkDestinationIsDir creates the directory dirPath if not exists
// on the target pod container
func checkDestinationIsDir(p *pods, podName, containerName string, dirPath string) error {
	logs.Info("Testing whether '%s/%s/%s:%s' is a directory.", p.namespace, podName, containerName, dirPath)

	cmd := []string{"test", "-d", dirPath}
	_, err := p.execCommand(podName, containerName, nil, cmd)

	return err
}

// makeTar tars the files and subdirectories of srcDir into tarDestDir (root directory within the tar file)
// than writes the created tar file to writer
func makeTar(srcPath, destPath string, writer io.Writer) error {
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

// recursiveTar tars recursively the content of srcDirPath
func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
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
