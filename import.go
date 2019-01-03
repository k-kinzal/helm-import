package main

import (
	"fmt"
	"net/url"
	"strings"
)

func Import(url url.URL) error {
	if strings.HasSuffix(url.Path, ".tgz") {
		return HttpImport(url)

	}
	if url.Host == "github.com" {
		return GithubImport(url)
	}

	return fmt.Errorf("is not match")
}
