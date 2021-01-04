package service

import (
	"fmt"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"
)

const (
	registryAPIVersion        = "v2"
	registryAPIBaseURL        = "%s/" + registryAPIVersion
	registryCatalogURL        = registryAPIBaseURL + "/" + "_catalog?n=10000"
	registryTagListURL        = registryAPIBaseURL + "/%s/tags/list"
	registryManifestURL       = registryAPIBaseURL + "/%s/manifests/%s"
	registryManifestDigestURL = registryAPIBaseURL + "/%s/manifests/%s"
)

func getCustomHeader() http.Header {
	return http.Header{
		"Accept": []string{"application/vnd.docker.distribution.manifest.v2+json"},
	}
}

func requestAndUnmarshal(method, specifiedURL string, target interface{}, reqHeader http.Header) (r *http.Response, err error) {
	utils.RequestHandle(method, specifiedURL, func(req *http.Request) error {
		req.Header = reqHeader
		return nil
	}, nil, func(req *http.Request, resp *http.Response) error {
		if target != nil {
			return utils.UnmarshalToJSON(resp.Body, target)
		}
		r = resp
		return nil
	})
	return
}

func GetRegistryCatalog() (repoList model.RegistryRepo, err error) {
	_, err = requestAndUnmarshal("GET", fmt.Sprintf(registryCatalogURL, registryURL()), &repoList, nil)
	return repoList, err
}

func GetRegistryImageTags(imageName string) (repoWithTags model.RegistryTags, err error) {
	_, err = requestAndUnmarshal("GET", fmt.Sprintf(registryTagListURL, registryURL(), imageName), &repoWithTags, nil)
	return repoWithTags, err
}

func GetRegistryManifest1(imageName, tagID string) (manifest1 model.RegistryManifest1, err error) {
	_, err = requestAndUnmarshal("GET", fmt.Sprintf(registryManifestURL, registryURL(), imageName, tagID), &manifest1, nil)
	return manifest1, err
}

func GetRegistryManifest2(imageName, tagID string) (manifest2 model.RegistryManifest2, err error) {
	_, err = requestAndUnmarshal("GET", fmt.Sprintf(registryManifestURL, registryURL(), imageName, tagID), &manifest2, getCustomHeader())
	return manifest2, err
}

func GetRegistryImageDigest(imageName, tagID string) (digest string, err error) {
	resp, err := requestAndUnmarshal("HEAD", fmt.Sprintf(registryManifestURL, registryURL(), imageName, tagID), nil, getCustomHeader())
	if resp.Header != nil {
		return strings.Trim(resp.Header.Get("Etag"), `"`), nil
	}
	return
}

func DeleteRegistryImageWithETag(imageName, tagID, digest string) (err error) {
	_, err = requestAndUnmarshal("DELETE", fmt.Sprintf(registryManifestDigestURL, registryURL(), imageName, digest), nil, getCustomHeader())
	return
}
