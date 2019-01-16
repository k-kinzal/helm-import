package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/helm/pkg/repo"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/cavaliercoder/grab"
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
	client := grab.NewClient()
	req, _ := grab.NewRequest(filepath, url.String())

	fmt.Fprintf(os.Stderr, "Downloading %v...\n", req.URL())
	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Fprintf(os.Stderr, "  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			fmt.Fprintf(os.Stderr, "  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return nil, fmt.Errorf("Download failed: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Download saved to %v \n", resp.Filename)

	return &resp.Filename, nil
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
