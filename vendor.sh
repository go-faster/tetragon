#!/bin/bash

VERSION=$1

grep -rl cilium/tetragon --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i 's|cilium/tetragon|go-faster/tetragon|g'
grep -rl quay.io --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i 's|quay.io|ghcr.io|g'
sed -i 's|-t "go-faster|-t "ghcr.io/go-faster|g' Makefile
sed -i 's|push go-faster|push ghcr.io/go-faster|g' Makefile
grep -rl v0.0.0-00010101000000-000000000000 --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i "s|v0.0.0-00010101000000-000000000000|$1|g"
grep -rl v0.8.3 --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i "s|v0.8.3|$1|g"

make vendor

git add vendor api pkg
git commit -a -m "vendor: go-faster/tetragon $VERSION"

git tag $VERSION
make image image-operator VERSION=$VERSION DOCKER_IMAGE_TAG=$VERSION

docker ghcr.io/go-faster/tetragon:latest:$VERSION ghcr.io/go-faster/tetragon:latest
docker push ghcr.io/go-faster/tetragon:$VERSION
docker push ghcr.io/go-faster/tetragon:latest

docker tag  ghcr.io/go-faster/tetragon-operator:$VERSION ghcr.io/go-faster/tetragon-operator:latest
docker push ghcr.io/go-faster/tetragon-operator:$VERSION
docker push ghcr.io/go-faster/tetragon-operator:latest

git push fork $VERSION