#!/bin/bash

archs=(amd64 arm64)
os=(linux darwin windows)

for os in ${os[@]}; do
    for arch in ${archs[@]}; do
        echo "Building ddns-cloudflare-agent for ${os} ${arch}"
        pushd app
        env GOOS=${os} GOARCH=${arch} go build -o ../bin/ddns-cloudflare-agent_${os}_${arch}
        popd
    done
done
echo "Done"
