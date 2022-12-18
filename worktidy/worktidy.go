package main


import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"

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

	// Check workspace requires for indirect requires.
	requires := make(map[string][]*modfile.Require)
	for _, use := range goWork.Use {
		goMod, err := readParseSubMod(use.Path)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		requires[goMod.Module.Mod.Path] = goMod.Require
	}

	vcsCmd, err := getVCS()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tags, err := vcsCmd.Tags(".")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	vTag := regexp.MustCompile(`^([^/]+/)?v\d+(\.\d+){0,2}$`)
	for _, tag := range tags {
		fmt.Printf("%s => %v\n", tag, vTag.Match([]byte(tag)))
	}
}
