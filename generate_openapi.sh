#!/bin/bash -x

if [ "$1" == "" ]; then
  echo "You must specify a branch of https://gitlab.tandashi.de/GameBase/swagger-rest-api.git as the first argument"
  exit 1
fi

mkdir tmp || exit 1
(
  cd tmp
  git clone https://gitlab.tandashi.de/GameBase/swagger-rest-api.git
  (
    cd swagger-rest-api
    git checkout "$1"
  )

  cmd.exe /C "npm install @openapitools/openapi-generator-cli"
  cmd.exe /C "npx openapi-generator generate -i swagger-rest-api/yaml-resolved/gamebase-api-specification.yaml -g go-gin-server -o out"

  rsync -avh --delete out/api/ ../api/
)

rm -rdf tmp
