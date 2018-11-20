package getter

import (
	"bytes"
	"fmt"
)

// Getter is an interface to support GET to the specified URL.
type Getter interface {
	//Get file content by url string
	Get(url string) (*bytes.Buffer, error)
}

type AuthGetter interface {
	Getter
	SetCredentials(username, password string)
}

// Constructor is the function for every getter which creates a specific instance
// according to the configuration
type Constructor func(URL string, Cert, Key, CA []byte) (Getter, error)

// Provider represents any getter and the schemes that it supports.
//
// For example, an HTTP provider may provide one getter that handles both
// 'http' and 'https' schemes.
type Provider struct {
	Schemes []string
	New     Constructor
}

// Provides returns true if the given scheme is supported by this Provider.
func (p Provider) Provides(scheme string) bool {
	for _, i := range p.Schemes {
		if i == scheme {
			return true
		}
	}
	return false
}

// Providers is a collection of Provider objects.
type Providers []Provider

// ByScheme returns a Provider that handles the given scheme.
//
// If no provider handles this scheme, this will return an error.
func (p Providers) ByScheme(scheme string) (Provider, error) {
	for _, pp := range p {
		if pp.Provides(scheme) {
			return pp, nil
		}
	}
	return Provider{}, fmt.Errorf("scheme %q not supported", scheme)
}

var all Providers

// ByScheme returns a getter for the given scheme.
//
// If the scheme is not supported, this will return an error.
func ByScheme(scheme string) (Provider, error) {
	return all.ByScheme(scheme)
}

func init() {
	all = Providers{
		{
			Schemes: []string{"http", "https"},
			New:     newHTTPGetter,
		},
	}
}
