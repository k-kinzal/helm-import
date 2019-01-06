package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/helm/pkg/repo"
	"net/http"
	"net/url"
	"os"
	"path"
)

func httpDownload(url url.URL) (*string, error) {
	cacheDir := HelmEnv.Home.Path("cache", "import", url.Host, path.Dir(url.Path))
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return nil, fmt.Errorf("%s: not create cache dir", cacheDir)
		}
	}
	filepath := path.Join(cacheDir, path.Base(url.Path))

	if _, err := os.Stat(filepath); err == nil {
		return &filepath, nil
	}

	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	return &filepath, nil
}

func HttpImport(url url.URL) error {
	filepath, err := httpDownload(url)
	if err != nil {
		return err
	}

	tmpdir, err := ioutil.TempDir("", "*")
	if err != nil {
		return err
	}

	if err := copy(*filepath, path.Join(tmpdir, path.Base(*filepath))); err != nil {
		return err
	}

	idx, err := repo.IndexDirectory(path.Join(tmpdir), Env.BaseUrl)
	if err != nil {
		return err
	}

	for _, versions := range idx.Entries {
		for _, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", version.Name, version.Version)
			if err := copy(*filepath, HelmEnv.Home.LocalRepository(filename)); err != nil {
				return err
			}
		}
	}


	index, err := repo.IndexDirectory(HelmEnv.Home.LocalRepository(), Env.BaseUrl)
	if err != nil {
		return err
	}

	return index.WriteFile(HelmEnv.Home.LocalRepository("index.yaml"), 0644)
}
