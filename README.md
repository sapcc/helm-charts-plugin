Helm charts plugin
------------------

Plugin to manage Helm charts in a directory/git repository.

## Install

```
helm plugin install https://github.com/sapcc/helm-charts --version=master
```

## Usage

```
Helm plugin to manage Helm charts in a directory.

Examples:
  $ helm charts list <path> <flags>

  flags:
    --exclude-dirs strings   List of (sub-)directories to exclude.
    --include-vendor         Also consider charts in the vendor folder.
    --only-path              Only output the chart path.
    --output-dir string      If given, results will be written to file in this directory.

  $ helm charts list-changed <path> <flags>

  flags:
    --exclude-dirs strings   List of (sub-)directories to exclude.
    --include-vendor         Also consider charts in the vendor folder.
    --only-path              Only output the chart path.
    --output-dir string      If given, results will be written to file in this directory.
    --remote string          The name of the git remote used to identify changes. (default "origin)"
    --branch string          The name of the branch used to identify changes. (default "master")
    --commit string          The commit used to identify changes. (default "HEAD")
```

## RELEASE

Releases are done via [goreleaser](https://github.com/goreleaser/goreleaser).  
Tag the new release, export the `GORELEASER_GITHUB_TOKEN` (needs `repo` scope) and run `make release`.
