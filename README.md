# Helm Import Plugin
[![Go Report Card](https://goreportcard.com/badge/github.com/k-kinzal/helm-import)](https://goreportcard.com/report/github.com/k-kinzal/helm-import)

Import Helm Chart where only public code into local repository.

## Get Started

```bash
$ helm plugin install https://github.com/k-kinzal/helm-import
```

## Usage

```bash
$ helm import URL
```

## Flags

```
-h, --help   help for import
```

## URL Patterns

* `https://example.com/path/to/file.tgz`
* `https://github.com/user/repo/`
* `https://github.com/user/repo/tree/branch`
* `https://github.com/user/repo/tree/branch/path/to`

helm-import supports these URL patterns.