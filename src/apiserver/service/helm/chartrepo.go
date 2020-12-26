package helm

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"gopkg.in/yaml.v2"

	"git/inspursoft/board/src/common/model"
)

const (
	helmName  = "helm"
	indexPath = "index.yaml"

	MuseumType  = 1
	DefaultType = 0
)

var supportedFiles = []string{"questions.yml", "questions.yaml"}

type ReleaseList struct {
	Next     string    `json:"Next,omitempty"`
	Releases []Release `json:"Releases,omitempty"`
}

type Release struct {
	Name      string `json:"Name,omitempty"`
	Namespace string `json:"Namespace,omitempty"`
	Revision  int32  `json:"Revision,omitempty"`
	Updated   string `json:"Updated,omitempty"`
	Status    string `json:"Status,omitempty"`
	Chart     string `json:"Chart,omitempty"`
}

type ReleaseStatus struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Info      *Info  `json:"info,omitempty"`
}

type Info struct {
	Status *Status `json:"status,omitempty"`
	// Description is human-friendly "log entry" about this release.
	Description string `json:"Description,omitempty"`
}

// Status defines the status of a release.
type Status struct {
	// Cluster resources as kubectl would print them.
	Resources string `json:"resources,omitempty"`
	// Contains the rendered templates/NOTES.txt if available
	Notes string `json:"notes,omitempty"`
}

// Entry represents a collection of parameters for chart repository
type Entry struct {
	Name     string
	URL      string
	Username string
	Password string
	Cert     []byte
	Key      []byte
	CA       []byte
	Type     int64
}

// ChartRepository represents a chart repository
type ChartRepository struct {
	Config    *Entry
	IndexFile *IndexFile
	Client    *HttpClient
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

func fillChartQuestions(chart *model.Chart) error {
	// generate the questions
	for _, file := range chart.Files {
		for _, f := range supportedFiles {
			if strings.EqualFold(f, file.Name) {
				var value model.Questions
				if err := yaml.Unmarshal([]byte(file.Contents), &value); err != nil {
					return err
				}
				chart.Questions = value.Questions
				return nil
			}
		}
	}
	return nil
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
	//fill the questions
	err = fillChartQuestions(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChartRepository) UploadChart(chartfile string) error {
	if r.Config.Type == MuseumType {
		absoluteChartURL, err := ResolveReferenceURL(r.Config.URL, "/api/charts")
		if err != nil {
			return fmt.Errorf("failed to make chart URL absolute: %v", err)
		}

		body, err := os.Open(chartfile)
		if err != nil {
			return err
		}
		defer body.Close()

		_, err = r.Client.Upload(absoluteChartURL, body)
		return err
	}
	return fmt.Errorf("the upload chart operation is not supported")
}

func (r *ChartRepository) DeleteChart(chartName, chartVersion string) error {
	if r.Config.Type == MuseumType {
		absoluteChartURL, err := ResolveReferenceURL(r.Config.URL, fmt.Sprintf("/api/charts/%s/%s", chartName, chartVersion))
		if err != nil {
			return fmt.Errorf("failed to make chart URL absolute: %v", err)
		}

		_, err = r.Client.Delete(absoluteChartURL)
		return err
	}
	return fmt.Errorf("the delete operation is not supported")

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

func (r *ChartRepository) Icon(versions ChartVersions) (string, string, error) {
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
	resp, err := client.Get(urlLocation)
	if err != nil {
		return "", "", fmt.Errorf("Error in HTTP GET of [%s], error: %s", urlLocation, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(urlLocation, "/")
	iconFilename := parts[len(parts)-1]
	iconData := base64.StdEncoding.EncodeToString(body)

	return iconData, iconFilename, nil
}

func (r *ChartRepository) InstallChart(chartName, chartVersion, releasename, namespace, values string, answers map[string]string, helmhost string) error {
	if helmhost == "" {
		return fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
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

	setValues := []string{}
	if answers != nil {
		for k, v := range answers {
			setValues = append(setValues, "--set", fmt.Sprintf("%s=%s", k, strings.Replace(v, ",", "\v", -1)))
		}
	}

	//create the release
	err = installChart(releasename, namespace, filepath.Join(targetdir, chartName), helmhost, setValues...)
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

	client, err := NewHTTPClient(cfg.URL, cfg.Username, cfg.Password, cfg.Cert, cfg.Key, cfg.CA)
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
func FindChartInRepoURL(repoURL, chartName, chartVersion string, cert, key, ca []byte, repotype int64) (string, error) {
	return FindChartInAuthRepoURL(repoURL, "", "", chartName, chartVersion, cert, key, ca, repotype)
}

// FindChartInAuthRepoURL finds chart in chart repository pointed by repoURL
// without adding repo to repositories, like FindChartInRepoURL,
// but it also receives credentials for the chart repository.
func FindChartInAuthRepoURL(repoURL, username, password, chartName, chartVersion string, cert, key, ca []byte, repotype int64) (string, error) {
	c := Entry{
		URL:      repoURL,
		Username: username,
		Password: password,
		Cert:     cert,
		Key:      key,
		CA:       ca,
		Type:     repotype,
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

func installChart(name, namespace, rootDir, helmhost string, setValues ...string) error {
	commands := make([]string, 0)
	commands = append([]string{"upgrade", "--install", "--namespace", namespace, name}, setValues...)
	commands = append(commands, rootDir)
	_, err := execHelmCommand(helmhost, commands...)
	if err != nil {
		//tiller maybe store the release, so we need to delete it manually.
		execHelmCommand(helmhost, "delete", "--purge", name)
	}
	return err
}

func DeleteRelease(release, helmhost string) error {
	combinedOutput, err := execHelmCommand(helmhost, "delete", "--purge", release)
	if err == nil || (err != nil && strings.Contains(combinedOutput, fmt.Sprintf("Error: release: \"%s\" not found", release))) {
		return nil
	}
	return errors.New(combinedOutput)
}

func ListAllReleases(helmhost string) (*ReleaseList, error) {
	combinedOutput, err := execHelmCommand(helmhost, "ls", "-a", "-m", fmt.Sprintf("%d", math.MaxInt32), "--output", "json")
	if err != nil {
		return nil, err
	}
	list := new(ReleaseList)
	if combinedOutput != "" {
		err = json.Unmarshal([]byte(combinedOutput), list)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func ListDeployedReleasesByNamespace(helmhost, namespace string) (*ReleaseList, error) {
	combinedOutput, err := execHelmCommand(helmhost, "ls", "--deployed", "-m", fmt.Sprintf("%d", math.MaxInt32), "--namespace", namespace, "--output", "json")
	if err != nil {
		return nil, err
	}
	list := new(ReleaseList)
	if combinedOutput != "" {
		err = json.Unmarshal([]byte(combinedOutput), list)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func GetRelease(release, helmhost string) (*Release, error) {
	combinedOutput, err := execHelmCommand(helmhost, "ls", "-m", "1", "-o", release, "--output", "json")
	if err != nil {
		return nil, err
	}
	list := new(ReleaseList)
	err = json.Unmarshal([]byte(combinedOutput), list)
	if err != nil {
		return nil, err
	}
	for i := range list.Releases {
		if release == list.Releases[i].Name {
			return &list.Releases[i], nil
		}
	}
	return nil, fmt.Errorf("can't find the release %s", release)
}

func GetReleaseValues(release, helmhost string) (string, error) {
	return execHelmCommand(helmhost, "get", "values", release)
}

func GetReleaseNotes(release, helmhost string) (string, error) {
	return execHelmCommand(helmhost, "get", "notes", release)
}

func GetReleaseManifest(release, helmhost string) (string, error) {
	return execHelmCommand(helmhost, "get", "manifest", release)
}

func GetReleaseStatus(release, helmhost string) (string, error) {
	combinedOutput, err := execHelmCommand(helmhost, "status", release, "--output", "json")
	if err != nil {
		return "", err
	}
	status := new(ReleaseStatus)
	err = json.Unmarshal([]byte(combinedOutput), status)
	if err != nil {
		return "", err
	}
	if status.Info != nil && status.Info.Status != nil {
		return status.Info.Status.Resources, nil
	}
	logs.Warning("can't get the release %s status", release)
	return "", nil
}

func execHelmCommand(helmhost string, args ...string) (string, error) {
	if helmhost == "" {
		return "", fmt.Errorf("You must specify the HELM_HOST environment when the apiserver starts")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, helmName, args...)
	cmd.Env = []string{fmt.Sprintf("%s=%s", "HELM_HOST", helmhost)}
	logs.Info("Execute command 'helm %s' with HELM_HOST=%s", strings.Join(args, " "), helmhost)
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		if context.DeadlineExceeded == ctx.Err() {
			logs.Error("Execute command 'helm %s' timeout, make sure the helm host %s is available", strings.Join(args, " "), helmhost)
		} else {
			logs.Error("Execute command 'helm %s' error: %+v", strings.Join(args, " "), err)
		}
		return "", errors.New(fmt.Sprintf("%s: %v", combinedOutput, err))
	}
	return string(combinedOutput), nil
}
