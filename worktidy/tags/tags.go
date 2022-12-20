package tags

import (
	"fmt"
	"regexp"
	"os"

	"golang.org/x/mod/module"
	"golang.org/x/tools/go/vcs"
)

type (
	TagReader struct {
		tags []string
	}
)

const (
	vRegexp = `v\d+(\.\d+){0,2}`
)

var (
	vcsList = []string{
		"hg",
		"git",
		"svn",
		"bzr",
	}
)

func NewTagReader() *TagReader {
	for _, v := range vcsList {
		if _, err := os.Stat(fmt.Sprintf(".%s", v)); err == nil {
			cmd :=  vcs.ByCmd(v)
			if tags, err := cmd.Tags("."); err == nil {
				return &TagReader{ tags: tags }
			}
		}
	}
	return &TagReader{ tags: make([]string, 0) }
}


func (p *TagReader) LatestFor(usePath string) string {
	if len(p.tags) == 0 {
		return ""
	}

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
	for _, tag := range p.tags {
		if !vTag.Match([]byte(tag)) {
			continue
		}
		vTags = append(vTags,
			module.Version{
				Path: "",
				Version: tag[len(prefix):],
			})
		}

	var latest string = ""
	if len(vTags) > 0 {
		// TODO: Sort is too much, we just need latest tag.
		module.Sort(vTags)
		latest = vTags[len(vTags)-1].Version
	}
	return latest
}
