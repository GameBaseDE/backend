#!/bin/bash -x

mkdir tmp || exit 1
(
  cd tmp
  git clone https://gitlab.tandashi.de/GameBase/swagger-rest-api.git || exit 1
  (
    cd swagger-rest-api
    if [ -z "$BRANCH" ]; then
      BRANCH="$1"
    fi
    if [ -z "$1" ]; then
      BRANCH=$(git branch --remote --sort=-committerdate | head -n1)
    else
      BRANCH="$1"
    fi
    git checkout $BRANCH || exit 1
  )

  npm install @openapitools/openapi-generator-cli afc11hn/gofmt.js || exit 1
  node_modules/@openapitools/openapi-generator-cli/bin/openapi-generator generate -i swagger-rest-api/yaml-resolved/gamebase-api-specification.yaml -g go-gin-server -o out || exit 1
)

node tools/rewrite.js || exit 1

# these files can't be moved into other directories because then go would consider them to be in another package
# which doesn't work since go doesn't support cyclic dependencies
mv openapi/kubernetes.go tmp/out/go
mv openapi/request_handlers.go tmp/out/go
mv openapi/http_* tmp/out/go

rsync --delete -achv tmp/out/api/ api/ || exit 1
rsync --delete -achv tmp/out/go/ openapi/ || exit 1

rm -rdf tmp
