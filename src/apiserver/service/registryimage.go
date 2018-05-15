package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"strings"
)

const (
	registryAPIVersion        = "v2"
	registryAPIBaseURL        = "%s/" + registryAPIVersion
	registryCatalogURL        = registryAPIBaseURL + "/" + "_catalog"
	registryTagListURL        = registryAPIBaseURL + "/%s/tags/list"
	registryManifestURL       = registryAPIBaseURL + "/%s/manifests/%s"
	registryManifestDigestURL = registryManifestURL + "/%s"
)

func requestAndUnmarshal(method, specifiedURL string, target interface{}, reqHeader map[string]string) (resp *http.Response, err error) {
	resp, err = utils.RequestHandle(method, specifiedURL, func(req *http.Request) error {
		for key, val := range reqHeader {
			req.Header.Set(key, val)
		}
		return nil
	}, nil)
	if err != nil {
		return
	}
	if target != nil {
		err = utils.UnmarshalToJSON(resp.Body, target)
		if err != nil {
			return
		}
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return
}

func getCustomHeader() (header map[string]string) {
	header = make(map[string]string)
	header["Accept"] = "application/vnd.docker.distribution.manifest.v2+json"
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
	_, err = requestAndUnmarshal("DELETE", fmt.Sprintf(registryManifestDigestURL, registryURL(), imageName, tagID, digest), nil, getCustomHeader())
	return
}
