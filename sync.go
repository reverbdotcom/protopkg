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

	"github.com/jtacoma/uritemplates"
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
			walker := &GithubWalker{
				Repo:   strings.Join(parts[0:2], "/"),
				Target: cfg.Path,
				Ref:    cfg.Ref,
			}

			err := walker.Visit(parts[2])

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

type GithubWalker struct {
	Repo   string
	Ref    string
	Target string
}

func (w *GithubWalker) Visit(src string) error {
	tpl, err := uritemplates.Parse("https://api.github.com/repos/{+repo}/contents/{+file}{?ref}")
	if err != nil {
		return err
	}

	path, err := tpl.Expand(map[string]interface{}{
		"repo": w.Repo,
		"file": src,
		"ref":  w.Ref,
	})
	if err != nil {
		return err
	}

	res, err := callGitHub(path)
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
			err := w.Visit(file.Path)
			if err != nil {
				return err
			}
		}

		if file.Type == "file" {
			err := w.downloadFile(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *GithubWalker) downloadFile(file GHFile) error {
	resp, err := callGitHub(file.URL)
	if err != nil {
		return err
	}

	proto := path.Join(w.Target, file.Path)
	os.MkdirAll(filepath.Dir(proto), os.ModePerm)
	f, err := os.Create(proto)

	if err != nil {
		return fmt.Errorf("could not create - %s", err)
	}

	io.Copy(f, resp)
	f.Close()

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

func callGitHub(url string) (io.Reader, error) {
	log.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if os.Getenv("GITHUB_TOKEN") != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	}

	resp, err := http.DefaultClient.Do(req)
	return resp.Body, err
}