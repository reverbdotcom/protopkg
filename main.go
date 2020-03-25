package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type ProtoDep struct {
	Path  string `json:"path"`
	Ref   string `json:"ref"`
	Local string `json:"local"`
}

type Manifest struct {
	Deps map[string]ProtoDep `json:"protos"`
}

func parseManifest() (*Manifest, error) {
	raw, err := ioutil.ReadFile("./protopkg.json")
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = json.Unmarshal(raw, &manifest)

	if err != nil {
		return nil, err
	}

	return &manifest, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "protopkg"
	app.Usage = "package manager for protocol buffers"

	app.Commands = []cli.Command{
		{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "set your GitHub token for pulling from private repositories",
			Action: func(c *cli.Context) error {
				return setToken(c.Args().First())
			},
		},
		{
			Name:    "local",
			Aliases: []string{"l"},
			Usage:   "sync a dependency based on the configured local path",
			Action: func(c *cli.Context) error {
				manifest, err := parseManifest()
				if err != nil {
					return err
				}

				return local(manifest, c.Args().First())
			},
		},
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "pull down the protos - protopkg sync",
			Action: func(c *cli.Context) error {
				manifest, err := parseManifest()
				if err != nil {
					return err
				}

				return sync(manifest)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "creates a new protopkg.json in the current directory - protopkg init",
			Action: func(c *cli.Context) error {
				contents := `{
					"protos": {
					}
				}`

				err := ioutil.WriteFile("./protopkg.json", []byte(contents), 0644)
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds a new proto dependency - protopkg add google/protos@HEAD ./protos/google",
			Action: func(c *cli.Context) error {
				var manifest Manifest

				if _, err := os.Stat("./protopkg.json"); os.IsNotExist(err) {
					manifest = Manifest{Deps: map[string]ProtoDep{}}
				} else {
					raw, err := ioutil.ReadFile("./protopkg.json")
					if err != nil {
						return err
					}

					err = json.Unmarshal(raw, &manifest)
					if err != nil {
						return err
					}
				}

				depAndRef := strings.Split(c.Args().First(), "@")
				target := c.Args().Get(1)

				if len(depAndRef) == 2 {
					manifest.Deps[depAndRef[0]] = ProtoDep{
						Path: target,
						Ref:  "HEAD",
					}
				} else {
					manifest.Deps[depAndRef[0]] = ProtoDep{
						Path: target,
					}
				}

				contents, err := json.MarshalIndent(manifest, "", "  ")
				if err != nil {
					return err
				}

				return ioutil.WriteFile("./protopkg.json", contents, 0644)
			},
		},
	}

	app.Run(os.Args)
}
