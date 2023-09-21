package main

import (
	"fmt"
	"log"
	"strings"
	cp "github.com/otiai10/copy"
)

func local(manifest *Manifest, name string) error {
	var dep ProtoDep
	var found bool

	for n, d := range manifest.Deps {
		depName := strings.Join(strings.Split(n, "/")[0:2], "/")
		log.Printf("checking %s against %s", depName, name)

		if name == depName {
			dep = d
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find package by the name - %s", name)
	}

	if dep.Local == "" {
		return fmt.Errorf("%s is not configured with a local path. Add local: 'your-path' to the protopkg.json", name)
	}

	log.Printf("dep %+v", dep)

	err := cp.Copy(dep.Local, dep.Path)
	if err != nil {
		return fmt.Errorf("could not perform a local sync of %s because %s", name, err)
	}

	return nil
}
