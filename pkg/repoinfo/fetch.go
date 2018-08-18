package repoinfo

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/pkg/errors"
)

func Fetch() (*Info, error) {
	r, err := packages.NewResolver(nil, packages.StdExcludeDirRegexps, logutils.NewStderrLog("getrepoinfo"))
	if err != nil {
		return nil, errors.Wrap(err, "can't make resolver")
	}

	prog, err := r.Resolve("./...")
	if err != nil {
		return nil, errors.Wrap(err, "can't resolve")
	}

	for _, p := range prog.Packages() {
		bp := p.BuildPackage()
		if bp.ImportComment == "" {
			continue
		}

		importComment := strings.TrimSuffix(bp.ImportComment, "/")
		if !strings.HasSuffix(importComment, p.Dir()) {
			return nil, fmt.Errorf("invalid import comment %q in dir %q", importComment, p.Dir())
		}

		repoCanonicalImportPath := strings.TrimSuffix(importComment, p.Dir())
		repoCanonicalImportPath = strings.TrimSuffix(repoCanonicalImportPath, "/")
		return &Info{
			CanonicalImportPath: repoCanonicalImportPath,
		}, nil
	}

	return &Info{}, nil
}
