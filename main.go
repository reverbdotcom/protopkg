package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

type ProtoDep struct {
	Path string `json:"path"`
	Ref  string `json:"ref"`
}

type Manifest struct {
	Deps map[string]ProtoDep `json:"protos"`
}

func main() {
	app := cli.NewApp()
	app.Name = "protopkg"
	app.Usage = "package manager for protocol buffers"

	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "pull down the protos",
			Action: func(c *cli.Context) error {
				raw, err := ioutil.ReadFile("./protopkg.json")
				if err != nil {
					return err
				}

				var manifest Manifest
				err = json.Unmarshal(raw, &manifest)

				if err != nil {
					return err
				}

				return sync(manifest)
			},
		},
	}

	app.Run(os.Args)
}
