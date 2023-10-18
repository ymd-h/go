package repo

import (
	"fmt"
	"regexp"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/mod/module"
)

type (
	VCS interface {
		Tags() ([]string, error)
		LatestRevision() (string, error)
		Timestamp(string) (string, error)
	}

	Git struct {
		dir string
	}

	LocalRepo struct {
		tags []string
		pseudoVersion string
	}
)

const (
	vRegexp = `v\d+(\.\d+){0,2}`
)

func (g *Git) Tags() ([]string, error) {
	cmd := exec.Command("git", "tag")
	cmd.Dir = g.dir
	cout, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("Fail to Get Tags with git: %w", err)
	}
	return strings.Split(string(cout), "\n"), nil
}

func (g *Git) LatestRevision() (string, error) {
	cmd := exec.Command("git", "log", "-n", "1", "--oneline", "--format='%H'")
	cmd.Dir = g.dir
	cout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Fail to Get Latest Revision with git: %w", err)
	}
	return string(cout[:12]), nil
}

func (g *Git) Timestamp(revision string) (string, error) {
	cmd := exec.Command("git", "log", "-n", "1", "--oneline", "--format='%cd'",
		"--date='unix'", revision)
	cmd.Dir = g.dir
	cout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Fail to Get Timestamp with git: %w", err)
	}

	unix, err := strconv.ParseInt(string(cout), 10, 64)
	if err != nil {
		return "", fmt.Errorf("Fail to Convert Timestamp: %w", err)
	}

	t := time.Unix(unix, 0).In(time.UTC)
	return t.Format("20060102150405"), nil
}

func NewLocalRepo(vcs VCS) (*LocalRepo, error) {
	tags, err := vcs.Tags()
	if err != nil {
		return nil, fmt.Errorf("Fail to Create LocalRepo: %w", err)
	}

	rev, err := vcs.LatestRevision()
	if err != nil {
		return nil, fmt.Errorf("Fail to Create LocalRepo: %w", err)
	}

	timestamp, err := vcs.Timestamp(rev)
	if err != nil {
		return nil, fmt.Errorf("Fail to Create LocalRepo: %w", err)
	}

	return &LocalRepo{
		tags: tags,
		pseudoVersion: fmt.Sprintf("v0.0.0-%s-%s", timestamp, rev),
	}, nil
}


func (p *LocalRepo) Version(usePath string) string {
	var latest string = p.pseudoVersion

	if len(p.tags) == 0 {
		return latest
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

	if len(vTags) > 0 {
		// TODO: Sort is too much, we just need latest tag.
		module.Sort(vTags)
		latest = vTags[len(vTags)-1].Version
	}
	return latest
}
