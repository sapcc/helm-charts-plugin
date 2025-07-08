<!--
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

Helm charts plugin
------------------

Helm plugin to manage multiple charts in a directory layout or git repository.

## Install

```
helm plugin install https://github.com/sapcc/helm-charts-plugin --version=master
```

## Usage

```
Helm plugin to manage Helm charts in a directory.

Examples:
  $ helm charts list <path> <flags>

  flags:
    --exclude-dirs strings   List of (sub-)directories to exclude.
    --only-path              Only output the chart path.
    --output-dir string      If given, results will be written to file in this directory.

  $ helm charts list-changed <path> <flags>

  flags:
    --exclude-dirs strings   List of (sub-)directories to exclude.
    --only-path              Only output the chart path.
    --output-dir string      If given, results will be written to file in this directory.
    --remote string          The name of the git remote used to identify changes. (default "origin)"
    --branch string          The name of the branch used to identify changes. (default "master")
    --commit string          The commit used to identify changes. (default "HEAD")
```

## RELEASE

Releases are done via [goreleaser](https://github.com/goreleaser/goreleaser).
Tag the new release, export the `GORELEASER_GITHUB_TOKEN` (needs `repo` scope) and run `make release`.
