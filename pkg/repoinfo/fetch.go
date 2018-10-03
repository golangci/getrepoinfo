package repoinfo

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/pkg/errors"
)

func extractAlias(prog *packages.Program) (string, error) {
	for _, p := range prog.Packages() {
		bp := p.BuildPackage()
		if bp.ImportComment == "" {
			continue
		}

		importComment := strings.TrimSuffix(bp.ImportComment, "/")
		var repoCanonicalImportPath string
		if p.Dir() == "." {
			repoCanonicalImportPath = importComment
		} else {
			if !strings.HasSuffix(importComment, p.Dir()) {
				return "", fmt.Errorf("invalid import comment %q in dir %q", importComment, p.Dir())
			}
			repoCanonicalImportPath = strings.TrimSuffix(importComment, p.Dir())
			repoCanonicalImportPath = strings.TrimSuffix(repoCanonicalImportPath, "/")
		}

		return repoCanonicalImportPath, nil
	}

	return "", nil
}

func Fetch(repo string) (*Info, error) {
	r, err := packages.NewResolver(nil, packages.StdExcludeDirRegexps, logutils.NewStderrLog("getrepoinfo"))
	if err != nil {
		return nil, errors.Wrap(err, "can't make resolver")
	}

	prog, err := r.Resolve("./...")
	if err != nil {
		return nil, errors.Wrap(err, "can't resolve")
	}

	alias, err := extractAlias(prog)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract alias")
	}
	if alias != "" {
		return &Info{
			CanonicalImportPath: alias,
		}, nil
	}

	if strings.ToLower(repo) != repo {
		return nil, fmt.Errorf("must set lowercased repo")
	}

	for _, p := range prog.Packages() {
		bp := p.BuildPackage()
		for _, imp := range bp.Imports {
			impLower := strings.ToLower(imp)
			if imp == impLower {
				continue
			}

			if strings.HasPrefix(impLower, repo+"/") {
				return &Info{
					CanonicalImportPath: imp[:len(repo)],
				}, nil
			}
		}
	}

	return &Info{}, nil
}
