#!/bin/bash

VERSION=$1

grep -rl cilium/tetragon --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i 's|cilium/tetragon|go-faster/tetragon|g'
grep -rl v0.0.0-00010101000000-000000000000 --exclude-dir=vendor --exclude-dir=.git --exclude vendor.sh | xargs sed -i "s|v0.0.0-00010101000000-000000000000|$1|g"

make vendor

git add vendor api pkg
git commit -a -m "vendor: update cilium/tetragon to go-faster/tetragon"

git tag -a $VERSION -m "v$VERSION"
make image VERSION=$VERSION DOCKER_IMAGE_TAG=$VERSION