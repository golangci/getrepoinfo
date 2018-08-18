package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/golangci/getrepoinfo/pkg/repoinfo"
	"github.com/pkg/errors"
)

func main() {
	if err := printRepoInfo(); err != nil {
		log.Fatal(err)
	}
}

func printRepoInfo() error {
	info, err := repoinfo.Fetch()
	if err != nil {
		return errors.Wrap(err, "can't get repo info")
	}

	if err = json.NewEncoder(os.Stdout).Encode(info); err != nil {
		return errors.Wrap(err, "can't json marshal info")
	}

	return nil
}
