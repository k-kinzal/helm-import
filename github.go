package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/repo"
	"net/url"
	"os"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func GithubImport(url url.URL) error {
	s := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(s) < 2 { // user/repo
		return fmt.Errorf("%s: is not Github repository URL", url.String())
	}

	repoUrl := fmt.Sprintf("%s://%s/%s", url.Scheme, url.Host, strings.Join(s[0:2], "/"))
	branch := "master"
	if len(s) >= 4 { // /user/repo/tree/branch
		branch = s[3]
	}
	pathChart := "/"
	if len(s) >= 5 { //usr/repo/tree/branch/path/to
		pathChart = fmt.Sprintf("/%s", strings.Join(s[4:], "/"))
	}

	tmpdir, err := ioutil.TempDir("", "*")
	if err != nil {
		return err
	}

	repository, err := git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stderr,
	})
	if err != nil {
		return err
	}

	if _, err := repository.Branch(branch); err != nil {
		return err
	}

	ch, err := chartutil.LoadDir(path.Join(tmpdir, pathChart))
	if err != nil {
		return err
	}

	if _, err := chartutil.Save(ch, HelmEnv.Home.LocalRepository()); err != nil {
		return err
	}

	index, err := repo.IndexDirectory(HelmEnv.Home.LocalRepository(), Env.BaseUrl)
	if err != nil {
		return err
	}

	return index.WriteFile(HelmEnv.Home.LocalRepository("index.yaml"), 0644)
}
