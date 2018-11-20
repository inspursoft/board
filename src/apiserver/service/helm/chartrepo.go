package helm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"

	"git/inspursoft/board/src/apiserver/service/helm/getter"
	"git/inspursoft/board/src/common/model"
)

const indexPath = "index.yaml"

type TemplateHandler func(string) error

// Entry represents a collection of parameters for chart repository
type Entry struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Cert     []byte `json:"cert"`
	Key      []byte `json:"key"`
	CA       []byte `json:"ca"`
}

// ChartRepository represents a chart repository
type ChartRepository struct {
	Config    *Entry
	IndexFile *IndexFile
	Client    getter.Getter
}

// Load loads a directory of charts as if it were a repository.
//
// It requires the presence of an index.yaml file in the directory.
func (r *ChartRepository) Load() error {
	dirInfo, err := os.Stat(r.Config.Name)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("%q is not a directory", r.Config.Name)
	}

	i, err := loadIndexFromFile(filepath.Join(r.Config.Name, indexPath))
	if err != nil {
		return nil
	}
	r.IndexFile = i

	return nil
}

// DownloadIndexFile fetches the index from a repository.
func (r *ChartRepository) downloadIndexFile() ([]byte, error) {
	var indexURL string
	parsedURL, err := url.Parse(r.Config.URL)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") + "/" + indexPath

	indexURL = parsedURL.String()

	r.setCredentials()
	resp, err := r.Client.Get(indexURL)
	if err != nil {
		return nil, err
	}

	index, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	return index, nil
}

// If HttpGetter is used, this method sets the configured repository credentials on the HttpGetter.
func (r *ChartRepository) setCredentials() {
	if t, ok := r.Client.(getter.AuthGetter); ok {
		t.SetCredentials(r.Config.Username, r.Config.Password)
	}
}

// Index generates an index for the chart repository and writes an index.yaml file.
func (r *ChartRepository) Index() error {
	index, err := yaml.Marshal(r.IndexFile)
	if err != nil {
		return err
	}
	err = os.MkdirAll(r.Config.Name, 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(r.Config.Name, indexPath), index, 0644)
}

func (r *ChartRepository) FetchTgz(url string) (*model.Chart, error) {
	var chart model.Chart

	logs.Debug("Fetching file %s", url)
	resp, err := r.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error in HTTP GET of [%s], error: %s", url, err)
	}

	gzf, err := gzip.NewReader(resp)
	if err != nil {
		return nil, err
	}
	defer gzf.Close()

	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			fallthrough
		case tar.TypeRegA:
			name := header.Name
			contents, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return nil, err
			}
			paths := strings.SplitN(name, "/", 2)
			if len(paths) == 2 {
				name = paths[1]
			}
			if name == "values.yaml" {
				chart.Values = string(contents)
			} else if name == "Chart.yaml" {
				var meta model.ChartMetadata
				if err = yaml.Unmarshal(contents, &meta); err != nil {
					logs.Error("unmarshal chart's Chart.yaml %s error: %+v", url, err)
					return nil, err
				}
				chart.Metadata = &meta
			} else if strings.HasPrefix(name, "templates/") {
				chart.Templates = append(chart.Templates, &model.File{
					Name:     name,
					Contents: string(contents),
				})
			} else {
				chart.Files = append(chart.Files, &model.File{
					Name:     name,
					Contents: string(contents),
				})
			}
		}
	}

	return &chart, nil
}

func (r *ChartRepository) formatChartURL(chartName, chartVersion string) (string, error) {
	errMsg := fmt.Sprintf("chart %q", chartName)
	if chartVersion != "" {
		errMsg = fmt.Sprintf("%s version %q", errMsg, chartVersion)
	}
	cv, err := r.IndexFile.Get(chartName, chartVersion)
	if err != nil {
		return "", fmt.Errorf("%s not found in %s repository", errMsg, r.Config.URL)
	}

	if len(cv.URLs) == 0 {
		return "", fmt.Errorf("%s has no downloadable URLs", errMsg)
	}

	chartURL := cv.URLs[0]

	absoluteChartURL, err := ResolveReferenceURL(r.Config.URL, chartURL)
	if err != nil {
		return "", fmt.Errorf("failed to make chart URL absolute: %v", err)
	}
	return absoluteChartURL, nil
}

func (r *ChartRepository) FetchChart(chartName, chartVersion string) (*model.Chart, error) {
	absoluteChartURL, err := r.formatChartURL(chartName, chartVersion)
	if err != nil {
		return nil, err
	}
	c, err := r.FetchTgz(absoluteChartURL)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChartRepository) DownloadChart(chartName, chartVersion, targetdir string) error {
	absoluteChartURL, err := r.formatChartURL(chartName, chartVersion)
	if err != nil {
		return err
	}
	resp, err := r.Client.Get(absoluteChartURL)
	if err != nil {
		return err
	}
	err = os.MkdirAll(targetdir, 0755)
	if err != nil {
		return err
	}
	gzf, err := gzip.NewReader(resp)
	if err != nil {
		return err
	}
	defer gzf.Close()

	tarReader := tar.NewReader(gzf)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath.Join(targetdir, header.Name), 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			fallthrough
		case tar.TypeRegA:
			name := header.Name
			contents, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return err
			}
			targetfile := filepath.Join(targetdir, name)
			err = os.MkdirAll(filepath.Dir(targetfile), 0755)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(targetfile, contents, 0777)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//func iconFromFile(versions ChartVersions) (string, string, error) {
//	for _, version := range versions {
//		if version.Dir == "" || version.Icon == "" {
//			continue
//		}
//
//		filename := filepath.Base(version.Icon)
//		iconFile := filepath.Join(filepath.Dir(version.Dir), filename)
//
//		bytes, err := ioutil.ReadFile(iconFile)
//		if err == nil {
//			return base64.StdEncoding.EncodeToString(bytes), filename, nil
//		}
//	}
//
//	return "", "", os.ErrNotExist
//}

func (r *ChartRepository) Icon(versions model.ChartVersions) (string, string, error) {
	//	data, file, err := iconFromFile(versions)
	//	if err == nil {
	//		return data, file, nil
	//	}

	if len(versions) == 0 || versions[0].Icon == "" {
		return "", "", nil
	}

	client := http.Client{
		Timeout: time.Second * 5,
	}
	urlLocation := versions[0].Icon
	iconUrl, err := url.Parse(urlLocation)
	if err != nil {
		return "", "", err
	}
	if iconUrl.Scheme == "https" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	resp, err := r.Client.Get(urlLocation)
	if err != nil {
		return "", "", fmt.Errorf("Error in HTTP GET of [%s], error: %s", urlLocation, err)
	}

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(urlLocation, "/")
	iconFilename := parts[len(parts)-1]
	iconData := base64.StdEncoding.EncodeToString(body)

	return iconData, iconFilename, nil
}

func (r *ChartRepository) InstallChart(chartName, chartVersion, releasename, namespace, values, helmhost string, handler TemplateHandler) error {
	targetdir, err := ioutil.TempDir("", "template")
	if err != nil {
		return err
	}
	defer os.RemoveAll(targetdir)

	err = r.DownloadChart(chartName, chartVersion, targetdir)
	if err != nil {
		return err
	}
	if values != "" {
		//override the values.yaml
		err = ioutil.WriteFile(filepath.Join(targetdir, chartName, "values.yaml"), []byte(values), 0777)
		if err != nil {
			return err
		}
	}

	//template the chart for resolve kubernetes elements
	templateInfo, err := templateChart(releasename, namespace, filepath.Join(targetdir, chartName), helmhost)
	if err != nil {
		return err
	}
	err = handler(templateInfo)
	if err != nil {
		return err
	}

	//create the release
	err = installChart(releasename, namespace, filepath.Join(targetdir, chartName), helmhost)
	if err != nil {
		return err
	}
	return nil
}

// NewChartRepository constructs ChartRepository
func NewChartRepository(cfg *Entry) (*ChartRepository, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid chart URL format: %s", cfg.URL)
	}

	get, err := getter.ByScheme(u.Scheme)
	if err != nil {
		return nil, fmt.Errorf("Could not find protocol handler for: %s", u.Scheme)
	}
	getterConstructor := get.New
	client, err := getterConstructor(cfg.URL, cfg.Cert, cfg.Key, cfg.CA)
	if err != nil {
		return nil, fmt.Errorf("Could not construct protocol handler for: %s error: %v", u.Scheme, err)
	}

	repo := &ChartRepository{
		Config:    cfg,
		IndexFile: NewIndexFile(),
		Client:    client,
	}
	indexBytes, err := repo.downloadIndexFile()
	if err != nil {
		return nil, fmt.Errorf("Looks like %q is not a valid chart repository or cannot be reached: %s", cfg.URL, err)
	}

	// Read the index file for the repository to get chart information and return chart URL
	repoIndex, err := loadIndex(indexBytes)
	if err != nil {
		return nil, err
	}

	repo.IndexFile = repoIndex
	return repo, nil
}

// FindChartInRepoURL finds chart in chart repository pointed by repoURL
// without adding repo to repositories
func FindChartInRepoURL(repoURL, chartName, chartVersion string, cert, key, ca []byte) (string, error) {
	return FindChartInAuthRepoURL(repoURL, "", "", chartName, chartVersion, cert, key, ca)
}

// FindChartInAuthRepoURL finds chart in chart repository pointed by repoURL
// without adding repo to repositories, like FindChartInRepoURL,
// but it also receives credentials for the chart repository.
func FindChartInAuthRepoURL(repoURL, username, password, chartName, chartVersion string, cert, key, ca []byte) (string, error) {
	c := Entry{
		URL:      repoURL,
		Username: username,
		Password: password,
		Cert:     cert,
		Key:      key,
		CA:       ca,
	}
	r, err := NewChartRepository(&c)
	if err != nil {
		return "", err
	}

	errMsg := fmt.Sprintf("chart %q", chartName)
	if chartVersion != "" {
		errMsg = fmt.Sprintf("%s version %q", errMsg, chartVersion)
	}
	cv, err := r.IndexFile.Get(chartName, chartVersion)
	if err != nil {
		return "", fmt.Errorf("%s not found in %s repository", errMsg, repoURL)
	}

	if len(cv.URLs) == 0 {
		return "", fmt.Errorf("%s has no downloadable URLs", errMsg)
	}

	chartURL := cv.URLs[0]

	absoluteChartURL, err := ResolveReferenceURL(repoURL, chartURL)
	if err != nil {
		return "", fmt.Errorf("failed to make chart URL absolute: %v", err)
	}

	return absoluteChartURL, nil
}

// ResolveReferenceURL resolves refURL relative to baseURL.
// If refURL is absolute, it simply returns refURL.
func ResolveReferenceURL(baseURL, refURL string) (string, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s as URL: %v", baseURL, err)
	}

	parsedRefURL, err := url.Parse(refURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s as URL: %v", refURL, err)
	}

	// if the base URL contains query string parameters,
	// propagate them to the child URL but only if the
	// refURL is relative to baseURL
	resolvedURL := parsedBaseURL.ResolveReference(parsedRefURL)
	if (resolvedURL.Hostname() == parsedBaseURL.Hostname()) && (resolvedURL.Port() == parsedBaseURL.Port()) {
		resolvedURL.RawQuery = parsedBaseURL.RawQuery
	}

	return resolvedURL.String(), nil
}

func loadIndexFromFile(indexPath string) (*IndexFile, error) {
	body, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	return loadIndex(body)
}

// loadIndex loads an index file and does minimal validity checking.
//
// This will fail if API Version is not set (ErrNoAPIVersion) or if the unmarshal fails.
func loadIndex(data []byte) (*IndexFile, error) {
	i := &IndexFile{}
	if err := yaml.Unmarshal(data, i); err != nil {
		return i, err
	}
	i.SortEntries()
	return i, nil
}

func templateChart(name, namespace, rootDir, helmhost string) (string, error) {
	commands := make([]string, 0)
	commands = append([]string{"template", "--namespace", namespace, "--name", name})
	commands = append(commands, rootDir)

	cmd := exec.Command(helmName, commands...)
	cmd.Env = []string{fmt.Sprintf("%s=%s", "HELM_HOST", helmhost)}
	stderrBuf := &bytes.Buffer{}
	stdoutBuf := &bytes.Buffer{}
	cmd.Stdout = stdoutBuf
	cmd.Stderr = stderrBuf
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to install chart %s. %s, error:%s", name, stderrBuf.String(), err.Error())
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("failed to install chart %s. %s, error:%s", name, stderrBuf.String(), err.Error())
	}
	return stdoutBuf.String(), nil
}

func installChart(name, namespace, rootDir, helmhost string) error {
	commands := make([]string, 0)
	commands = append([]string{"upgrade", "--install", "--namespace", namespace, name})
	commands = append(commands, rootDir)

	cmd := exec.Command(helmName, commands...)
	cmd.Env = []string{fmt.Sprintf("%s=%s", "HELM_HOST", helmhost)}
	stderrBuf := &bytes.Buffer{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = stderrBuf
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to install chart %s. %s, error:%s", name, stderrBuf.String(), err.Error())
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to install chart %s. %s, error:%s", name, stderrBuf.String(), err.Error())
	}
	return nil
}

func deleteRelease(release, helmhost string) error {
	cmd := exec.Command(helmName, "delete", "--purge", release)
	cmd.Env = []string{fmt.Sprintf("%s=%s", "HELM_HOST", helmhost)}
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil && combinedOutput != nil && strings.Contains(string(combinedOutput), fmt.Sprintf("Error: release: \"%s\" not found", release)) {
		return nil
	}
	return errors.New(string(combinedOutput))
}
