package main


import (
	"fmt"
	"go/token"
	"go/parser"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/ymd-h/go/worktidy/tags"
)

type (
	Module struct {
		Version string
		Require []*modfile.Require
	}
)


func readParseWork() (*modfile.WorkFile, error) {
	workByte, err := os.ReadFile("go.work")
	if err != nil {
		return nil, err
	}
	return modfile.ParseWork("go.work", workByte, nil)
}


func readParseSubMod(path string) (*modfile.File, error) {
	modByte, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, err
	}
	return modfile.Parse("go.mod", modByte, nil)
}


func main() {
	goWork, err := readParseWork()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Ensure update non-workspace requires by `go mod tidy -e`.
	var wg sync.WaitGroup
	for _, use := range goWork.Use {
		wg.Add(1)
		go func(path string){
			defer wg.Done()
			// Since `go mod tidy` is implemented at internal package,
			// we call external command instead.
			// Workspace modules might not be published yet, so that
			// we ignore missing error by adding `-e` option.
			cmd := exec.Command("go", "mod", "tidy", "-e")
			abs, _ := filepath.Abs(path)
			cmd.Dir = abs
			cout, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%s\n%s\n", cout, err.Error())
			}
		}(use.Path)
	}
	wg.Wait()

	// Check workspace requires and latest versions.
	mod := make(map[string]*Module)
	tag := tags.NewTagReader()
	for _, use := range goWork.Use {
		usePath := use.Path
		goMod, err := readParseSubMod(usePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		modPath := goMod.Module.Mod.Path
		latest := tag.LatestFor(usePath)
		fmt.Printf("%s: %s\n", modPath, latest)

		mod[modPath] = &Module{
			Version: latest,
			Require: goMod.Require,
		}
	}

	for _, use := range goWork.Use {
		fset := token.NewFileSet()
		astMap, err := parser.ParseDir(
			fset,
			filepath.Clean(use.Path),
			func(f fs.FileInfo) bool {
				return !strings.HasSuffix(f.Name(), "_test.go")
			},
			parser.ImportsOnly,
		)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		req := make(map[string]struct{})
		for _, pkg := range astMap {
			for _, astFile := range pkg.Files {
				for _, i := range astFile.Imports {
					if i.Path != nil {
						v := i.Path.Value[1:len(i.Path.Value)-1]
						for modPath, _ := range mod {
							if strings.HasPrefix(v, modPath) {
								req[modPath] = struct{}{}
								break
							}
						}
					}
				}
			}
		}
		fmt.Printf("Submodule: %s\n", use.Path)
		for p, _ := range req {
			fmt.Printf("  Depends on %s\n", p)
		}
	}
}
