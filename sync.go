package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		fmt.Printf("♻️  syncing %s to %s\n", path, cfg.Path)

		go func(p string, c ProtoDep) {
			parts := strings.SplitN(p, "/", 3)

			err := Walk(
				parts[2],
				strings.Join(parts[0:2], "/"),
				c.Path,
				c.Ref,
			)

			if err != nil {
				fmt.Printf("👎 error! - %s\n", err)
			}
			depC <- true
		}(path, cfg)
	}

	for range manifest.Deps {
		<-depC
	}

	fmt.Printf("📦 finished!\n")
	return nil
}

type GithubWalker struct {
	Repo   string
	Ref    string
	Target string
	Base   string
}

func Walk(source, repo, target, ref string) error {
	walker := &GithubWalker{
		Repo:   repo,
		Target: target,
		Ref:    ref,
		Base:   source,
	}

	return walker.Visit(source)
}

func (w *GithubWalker) Visit(src string) error {
	tpl, err := uritemplates.Parse("https://api.github.com/repos/{+repo}/contents/{+file}{?ref}")
	if err != nil {
		return err
	}

	ref := w.Ref
	if ref == "" {
		ref = "HEAD"
	}

	params := map[string]interface{}{
		"repo": w.Repo,
		"file": src,
		"ref":  ref,
	}

	path, err := tpl.Expand(params)
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
			err := w.downloadFile(file, src)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *GithubWalker) downloadFile(file GHFile, base string) error {
	resp, err := callGitHub(file.URL)
	if err != nil {
		return err
	}

	var protoPath string
	if filepath.Ext(w.Target) == ".proto" {
		protoPath = w.Target
	} else {
		p, _ := filepath.Rel(w.Base, file.Path)
		protoPath = path.Join(w.Target, p)
	}

	os.MkdirAll(filepath.Dir(protoPath), os.ModePerm)
	f, err := os.Create(protoPath)

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
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	if os.Getenv("GITHUB_TOKEN") != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
	}

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch %s", url)
	}

	return resp.Body, err
}
