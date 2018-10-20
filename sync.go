package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GHFile struct {
	Type string `json:"type"`
	Path string `json:"path"`
	URL  string `json:"download_url"`
}

// Sync will take the manifest and copy all of the containing proto files
// to the target.
func sync(manifest Manifest) error {
	depC := make(chan bool, len(manifest.Deps))

	for path, cfg := range manifest.Deps {
		log.Printf("syncing %s to %s", path, cfg.Path)
		parts := strings.SplitN(path, "/", 3)

		go func() {
			err := visit(parts[2], strings.Join(parts[0:2], "/"), cfg.Path)
			if err != nil {
				log.Printf("error! - %s", err)
			}
			depC <- true
		}()
	}

	for range manifest.Deps {
		<-depC
	}

	log.Printf("finished")
	return nil
}

func visit(src string, repo string, target string) error {
	res, err := callGitHub(fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, src))
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadAll(res)
	files, err := decodeFile(raw)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Type == "dir" {
			err := visit(file.Path, repo, target)
			if err != nil {
				return err
			}
		}

		if file.Type == "file" {
			err := downloadFile(file, target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func decodeFile(raw []byte) ([]GHFile, error) {
	var files []GHFile
	err := json.Unmarshal(raw, &files)

	if err != nil {
		var file GHFile
		err = json.Unmarshal(raw, &file)
		if err != nil {
			return files, err
		}

		files = append(files, file)
	}

	return files, err
}

func downloadFile(file GHFile, target string) error {
	resp, err := callGitHub(file.URL)
	if err != nil {
		return err
	}

	proto := path.Join(target, file.Path)
	os.MkdirAll(filepath.Dir(proto), os.ModePerm)
	f, err := os.Create(proto)

	if err != nil {
		return fmt.Errorf("could not create - %s", err)
	}

	io.Copy(f, resp)
	f.Close()

	return nil
}

func callGitHub(url string) (io.Reader, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if os.Getenv("GITHUB_TOKEN") != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	}

	resp, err := http.DefaultClient.Do(req)
	return resp.Body, err
}
