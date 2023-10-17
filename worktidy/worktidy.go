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

	"github.com/ymd-h/go/sets"
	"github.com/ymd-h/go/worktidy/repo"
	"github.com/ymd-h/go/worktidy/resolve"
)

type (
	Module struct {
		Version string
		ModFile *modfile.File
		UsePath string
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
	tag, err := repo.NewLocalRepo(&repo.Git{dir: "."})
	if err != nil {
		fmt.Printf("Fail to Get Local Repo: %w\n", err)
		return
	}

	for _, use := range goWork.Use {
		usePath := use.Path
		goMod, err := readParseSubMod(usePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		modPath := goMod.Module.Mod.Path
		latest := tag.Version(usePath)
		fmt.Printf("%s: %s\n", modPath, latest)

		mod[modPath] = &Module{
			Version: latest,
			ModFile: goMod,
			UsePath: usePath,
		}
	}

	directDeps := make(map[string]sets.ISet[string], len(goWork.Use))
	for p, m := range mod {
		req := sets.New[string]()
		fset := token.NewFileSet()

		err := filepath.WalkDir(filepath.Clean(m.UsePath),
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				if (!strings.HasSuffix(path, ".go")) ||
					strings.HasSuffix(path, "_test.go") {
					return nil
				}
				astF, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
				if err != nil {
					return err
				}
				for _, i := range astF.Imports {
					if i.Path == nil {
						continue
					}
					v := i.Path.Value[1:len(i.Path.Value)-1]
					for modPath, modV := range mod {
						if m.UsePath == modV.UsePath {
							// Skip Self
							continue
						}

						if strings.HasPrefix(v, modPath) {
							req.Add(modPath)
							break
						}
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		fmt.Printf("Submodule: %s\n", m.UsePath)
		for _, r := range req.ToSlice() {
			fmt.Printf("  Depends on %s\n", r)
		}

		directDeps[p] = req
	}

	totalDeps, err := resolve.Resolve(directDeps)
	if err != nil {
		fmt.Printf("Fail to Resolve Dependancy Graph: %w", err)
		return
	}

	for p, m := range mod {
		d := directDeps[p]
		t := totalDeps[p]

		// Indirect
		i := sets.New[string]()
		for _, td := range t.ToSlice() {
			if !d.Has(td) {
				i.Add(td)
			}
		}

		req := sets.New[string]()
		for _, r := range m.ModFile.Require {
			if r == nil {
				continue
			}
			req.Add(r.Mod.Path)
		}

		for _, dd := range d.ToSlice() {
			m.ModFile.AddNewRequire(dd, mod[dd].Version, false)

			for _, r := range mod[dd].ModFile.Require {
				if (r == nil) || (req.Has(r.Mod.Path)) {
					continue
				}
				req.Add(r.Mod.Path)
				m.ModFile.AddNewRequire(
					r.Mod.Path,
					r.Mod.Version,
					true,
				)
			}
		}

		for _, id := range i.ToSlice() {
			m.ModFile.AddNewRequire(id, mod[id].Version, true)

			for _, r := range mod[id].ModFile.Require {
				if (r == nil) || (req.Has(r.Mod.Path)) {
					continue
				}
				req.Add(r.Mod.Path)
				m.ModFile.AddNewRequire(
					r.Mod.Path,
					r.Mod.Version,
					true,
				)
			}
		}

		m.ModFile.SortBlocks()
		m.ModFile.Cleanup()
		b, err := m.ModFile.Format()
		if err != nil {
			fmt.Printf("Fail to Format %s: %w\n", m.UsePath, err)
			continue
		}

		err = os.WriteFile(filepath.Join(m.UsePath, "go.mod"), b, 0664)
		if err != nil {
			fmt.Printf("Fail to Write %s/go.mod: %w", m.UsePath, err)
			continue
		}
	}
}
