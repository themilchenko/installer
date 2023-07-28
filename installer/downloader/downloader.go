package downloader

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	DefaultDownloadPath = "/home/milchenko/programming/aktiv/installer/installer/tmp/"
)

func downloadFile(fileName, URL string) error {
	out, err := os.Create(DefaultDownloadPath + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Setting up https client
	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
	}
	tr := &http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client := http.Client{Transport: tr}

	data := url.Values{}
	data.Set("fileName", fileName)

	u, _ := url.ParseRequestURI("https://" + URL)
	u.Path = "/api/download"
	urlStr := u.String()

	r, _ := http.NewRequest(
		http.MethodGet,
		urlStr,
		nil,
	)
	r.URL.RawQuery = data.Encode()

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadAllFiles(fileNames []string, url string) {
	for _, fileName := range fileNames {
		err := downloadFile(fileName, url)
		if err != nil {
			log.Printf("Cannot load file with name %s", fileName)
		}
	}
}
