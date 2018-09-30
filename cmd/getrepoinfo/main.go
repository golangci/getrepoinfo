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
	var ret interface{}
	info, err := repoinfo.Fetch()
	if err != nil {
		ret = struct {
			Error string
		}{
			Error: err.Error(),
		}
	} else {
		ret = info
	}

	if err = json.NewEncoder(os.Stdout).Encode(ret); err != nil {
		return errors.Wrap(err, "can't json marshal")
	}

	return nil
}
