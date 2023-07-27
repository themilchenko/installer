package downloader

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultDownloadPath = "/home/milchenko/programming/aktiv/base.agent/installer/dump/"
)

func downloadFile(fileName, URL string) error {
	out, err := os.Create(defaultDownloadPath + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	data := url.Values{}
	data.Set("fileName", fileName)

	u, _ := url.ParseRequestURI("http://" + URL)
	u.Path = "/api/download"
	urlStr := u.String()

	client := &http.Client{}
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
