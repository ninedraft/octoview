package octoview

import (
	"encoding/json"
	"image"
	"net"
	"net/http"
	"time"

	"github.com/golang/freetype"
)

type PreviewGenerator struct {
	client *http.Client
}

func NewPreviewGenerator(params ...func(config *config)) PreviewGenerator {
	var config = &config{
		client: DefaultHTTPClient(),
	}
	for _, param := range params {
		param(config)
	}
	return PreviewGenerator{
		client: config.client,
	}
}

func (preview PreviewGenerator) getRepo(link string) (Repo, error) {
	var resp, errGetRepo = preview.client.Get(link)
	if errGetRepo != nil {
		return Repo{}, errGetRepo
	}
	defer resp.Body.Close()
	var repo Repo
	return repo, json.NewDecoder(resp.Body).Decode(&repo)
}

type metaAssets struct {
	ownerAvatar image.Image
	repoLogo    image.Image
}

func (preview PreviewGenerator) generateRaster(w, h uint, assets metaAssets, repo Repo) image.Image {
	var rect = image.Rect(0, 0, int(w), int(h))
	var img = image.NewRGBA(rect)
	var printer = freetype.NewContext()
	printer.SetDst(img)
	// TODO(ninedraft): add preview generation
	return img
}

func WithTransport(transport *http.Transport) func(config *config) {
	return func(config *config) {
		config.client.Transport = transport
	}
}

func WithTimeout(timeout time.Duration) func(config *config) {
	return func(config *config) {
		config.client.Timeout = timeout
	}
}

type config struct {
	client *http.Client
}

func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}
