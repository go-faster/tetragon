#!/bin/bash

VERSION=$1

if [[ $# -eq 0 ]] ; then
    echo 'Version is not set'
    exit 1
fi

grep -rl cilium/tetragon --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i 's|cilium/tetragon|go-faster/tetragon|g'
sed -i 's|quay.io|ghcr.io|g' install/kubernetes/values.yaml
sed -i 's|-t "go-faster|-t "ghcr.io/go-faster|g' Makefile
sed -i 's|push go-faster|push ghcr.io/go-faster|g' Makefile
grep -rl v0.0.0-00010101000000-000000000000 --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i "s|v0.0.0-00010101000000-000000000000|$1|g"
sed -i "s|0.8.3|${VERSION:1}|g" install/kubernetes/Chart.yaml
sed -i "s|v0.8.3|$1|g" install/kubernetes/values.yaml

rm -rf pkg/k8s/vendor
rm -f pkg/k8s/go.sum
rm -f pkg/k8s/go.mod

rm -rf api/vendor
rm -f api/go.sum
rm -f api/go.mod

sed -i "s|github.com/go-faster/tetragon/api => ./api||g" go.mod
sed -i "s|github.com/go-faster/tetragon/pkg/k8s => ./pkg/k8s||g" go.mod
sed -i "s|github.com/go-faster/tetragon/api $VERSION||g" go.mod
sed -i "s|github.com/go-faster/tetragon/pkg/k8s $VERSION||g" go.mod

make vendor

git add vendor api pkg go.mod go.sum
git commit -a -m "vendor: go-faster/tetragon $VERSION"

git tag $VERSION
make image image-operator VERSION=$VERSION DOCKER_IMAGE_TAG=$VERSION

docker push ghcr.io/go-faster/tetragon:$VERSION
docker push ghcr.io/go-faster/tetragon-operator:$VERSION

git push fork $VERSION
git push -f fork main
