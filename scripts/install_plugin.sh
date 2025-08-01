#!/bin/sh -e

# SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company
#
# SPDX-License-Identifier: Apache-2.0

# Copied w/ love from the excellent hypnoglow/helm-s3

if [ -n "${HELM_OUTDATED_DEPENDENCIES_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

version="$(grep "version" plugin.yaml | cut -d '"' -f 2)"
echo "Downloading and installing helm-charts-plugin v${version} ..."

url=""
if [ "$(uname)" = "Darwin" ]; then
    url="https://github.com/sapcc/helm-charts-plugin/releases/download/v${version}/helm-charts-plugin_${version}_darwin_amd64.tar.gz"
elif [ "$(uname)" = "Linux" ] ; then
    url="https://github.com/sapcc/helm-charts-plugin/releases/download/v${version}/helm-charts-plugin_${version}_linux_amd64.tar.gz"
else
    url="https://github.com/sapcc/helm-charts-plugin/releases/download/v${version}/helm-charts-plugin_${version}_windows_amd64.tar.gz"
fi

echo "$url"

mkdir -p "bin"
mkdir -p "releases/v${version}"

# Download with curl if possible.
if [ -x "$(command -v curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/v${version}.tar.gz"
else
    wget -q "${url}" -O "releases/v${version}.tar.gz"
fi
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"
mv "releases/v${version}/bin/helm-charts-plugin" "bin/helm-charts" || \
    mv "releases/v${version}/bin/helm-charts-plugin.exe" "bin/helm-charts"
