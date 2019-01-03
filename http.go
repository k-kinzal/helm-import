package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"k8s.io/helm/pkg/repo"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

func HttpDownload(url url.URL) (*string, error) {
	tmpdir, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile(tmpdir, "*.tgz")
	if err != nil {
		log.Fatal(err)
	}
	filepath := tmpfile.Name()

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
	filepath, err := HttpDownload(url)
	if err != nil {
		return err
	}

	i1, err := repo.LoadIndexFile(HelmEnv.Home.LocalRepository("index.yaml"))
	if err != nil {
		return err
	}

	i2, err := repo.IndexDirectory(path.Dir(*filepath), Env.BaseUrl)
	if err != nil {
		return err
	}
	for name, versions := range i2.Entries {
		for index, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", version.Name, version.Version)
			if err := copy(*filepath, HelmEnv.Home.LocalRepository(filename)); err != nil {
				return err
			}

			version.URLs[0] = fmt.Sprintf("%s/%s", Env.BaseUrl, filename)
			versions[index] = version
		}
		i2.Entries[name] = versions
	}

	i1.Merge(i2)

	return i1.WriteFile(HelmEnv.Home.LocalRepository("index.yaml"), 0644)
}
