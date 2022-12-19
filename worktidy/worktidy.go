package main


import (
	"errors"
	"fmt"
	"go/token"
	"go/parser"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/vcs"
)

var (
	vcsList = []string{
		"hg",
		"git",
		"svn",
		"bzr",
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


func getVCS() (*vcs.Cmd, error) {
	for _, v := range vcsList {
		if _, err := os.Stat(fmt.Sprintf(".%s", v)); err == nil {
			return vcs.ByCmd(v), nil
		}
	}
	return nil, errors.New("Unknown VCS")
}

func getTags() []string {
	vcsCmd, err := getVCS()
	if err != nil {
		return []string{}
	}

	tags, err := vcsCmd.Tags(".")
	if err != nil {
		return []string{}
	}

	return tags
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
			cmd := exec.Command("go", "mod", "tidy", "-e")
			cmd.Path = filepath.Clean(path)
			cout, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%s\n%s\n", cout, err.Error())
			}
		}(use.Path)
	}
	wg.Wait()

	// Check workspace requires and latest versions.
	requires := make(map[string][]*modfile.Require)
	versions := make(map[string]module.Version)
	tags := getTags()
	vRegexp := `v\d+(\.\d+){0,2}`
	for _, use := range goWork.Use {
		usePath := use.Path
		goMod, err := readParseSubMod(usePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		modPath := goMod.Module.Mod.Path
		requires[modPath] = goMod.Require

		var prefix string
		switch {
		case usePath == ".":
			prefix = ""
		case len(usePath) > 2 && usePath[:2] == "./":
			prefix = usePath[2:] + "/"
		default:
			prefix = usePath + "/"
		}

		vTag := regexp.MustCompile(fmt.Sprintf(`^%s%s$`, prefix, vRegexp))
		vTags := []module.Version{}
		for _, tag := range tags {
			if !vTag.Match([]byte(tag)) {
				continue
			}
			vTags = append(vTags,
				module.Version{
					Path: modPath,
					Version: tag[len(prefix):],
				})
		}
		if len(vTags) > 0 {
			module.Sort(vTags)
			versions[modPath] = vTags[len(vTags)-1]
		} else {
			versions[modPath] = module.Version{
				Path: modPath,
				Version: "",
			}
		}
	}


	for _, use := range goWork.Use {
		fmt.Printf("Check Import at %s\n", use.Path)
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
		for _, pkg := range astMap {
			fmt.Printf("pkg.Name: %s\n", pkg.Name)
			for pkgID, pkgObj := range pkg.Imports {
				fmt.Printf("pkgID: %s, pkgObj.Name: %s\n", pkgID, pkgObj.Name)
			}
		}
	}
}
