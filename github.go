package main

import (
	"fmt"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/repo"
	"net/url"
	"os"
	"path"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func gitDownload(url url.URL) (*string, error) {
	s := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(s) < 2 { // user/repo
		return nil, fmt.Errorf("%s: is not Github repository URL", url.String())
	}

	repoUrl, err := url.Parse(fmt.Sprintf("%s://%s/%s", url.Scheme, url.Host, strings.Join(s[0:2], "/")))
	if err != nil {
		return nil, err
	}

	branch := "master"
	if len(s) >= 4 { // /user/repo/tree/branch
		branch = s[3]
	}
	pathChart := "/"
	if len(s) >= 5 { //usr/repo/tree/branch/path/to
		pathChart = fmt.Sprintf("/%s", strings.Join(s[4:], "/"))
	}

	cacheDir := HelmEnv.Home.Path("cache", "import", repoUrl.Host, repoUrl.Path)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return nil, fmt.Errorf("%s: not create cache dir", cacheDir)
		}
	}

	repository, err := git.PlainOpen(cacheDir)
	if err != nil {
		repository, err = git.PlainClone(cacheDir, false, &git.CloneOptions{
			URL:      repoUrl.String(),
			Progress: os.Stderr,
		})
		if err != nil {
			return nil, err
		}
	}

	if _, err := repository.Branch(branch); err != nil {
		return nil, err
	}

	if err := repository.Fetch(&git.FetchOptions{Progress: os.Stderr, Force: true}); err != nil {
		if err.Error() != "already up-to-date" {
			return nil, err
		}
	}

	dirpath := path.Join(cacheDir, pathChart)

	return &dirpath, nil
}

func GithubImport(url url.URL) error {
	dirpath, err := gitDownload(url)
	if err != nil {
		return err
	}

	ch, err := chartutil.LoadDir(*dirpath)
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
