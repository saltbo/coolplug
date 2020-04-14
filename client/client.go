package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/saltbo/coolplug/model"
)

type Client interface {
	PluginList() ([]model.Plugin, error)
	PluginInstall(pi url.Values, filepath string) error
	PluginUninstall(id string) error
}

type HTTPClient struct {
	*resty.Client
}

func NewHTTPClient(hostURL string) *HTTPClient {
	client := resty.New()
	client.SetHostURL(hostURL)
	return &HTTPClient{Client: client}
}

func (hc *HTTPClient) PluginList() ([]model.Plugin, error) {
	mps := make([]model.Plugin, 0)
	if _, err := hc.R().SetResult(&mps).Get("/api/plugins"); err != nil {
		return nil, err
	}

	return mps, nil
}

func (hc *HTTPClient) PluginInstall(pi url.Values, filepath string) (err error) {
	resp, err := hc.R().SetFormDataFromValues(pi).SetFile("file", filepath).Post("/api/plugins")
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(resp.String())
	}

	return
}

func (hc *HTTPClient) PluginUninstall(id string) (err error) {
	_, err = hc.R().Delete("/api/plugins/" + id)
	return
}
