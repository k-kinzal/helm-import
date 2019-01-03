# Helm Import Plugin

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

* https://example.com/path/to/file.tgz
* helm import https://github.com/user/repo/
* helm import https://github.com/user/repo/tree/branch
* helm import https://github.com/user/repo/tree/branch/path/to

helm-import supports these URL patterns.