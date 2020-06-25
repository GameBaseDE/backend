[![go report card](https://goreportcard.com/badge/gitlab.tandashi.de/GameBase/gamebase-backend)](https://goreportcard.com/report/gitlab.tandashi.de/GameBase/gamebase-backend)
[![pipeline status](https://gitlab.tandashi.de/GameBase/gamebase-backend/badges/master/pipeline.svg)](https://gitlab.tandashi.de/GameBase/gamebase-backend/commits/master)
[![coverage report](https://gitlab.tandashi.de/GameBase/gamebase-backend/badges/master/coverage.svg)](https://gitlab.tandashi.de/GameBase/gamebase-backend/-/commits/master)

[![Bugs](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=bugs)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Code Smells](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=code_smells)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Coverage](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=coverage)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Duplicated Lines (%)](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=duplicated_lines_density)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Lines of Code](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=ncloc)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Maintainability Rating](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=sqale_rating)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Quality Gate Status](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=alert_status)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Reliability Rating](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=reliability_rating)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Security Rating](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=security_rating)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Technical Debt](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=sqale_index)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)
[![Vulnerabilities](https://sonarqube.gahr.dev/api/project_badges/measure?project=gamebase-daemon&metric=vulnerabilities)](https://sonarqube.gahr.dev/dashboard?id=gamebase-daemon)

# GameBase Backend
This is the backend service for the GameBase game server management platform.
GameBase makes it easy to manage game servers inside a Kubernetes cluster.

The backend is implemented as a REST API based on the 
[GameBase REST API swagger definition](https://gitlab.tandashi.de/GameBase/swagger-rest-api). 
It is supposed to be used in conjunction with the 
[GameBase Frontend](https://gitlab.tandashi.de/GameBase/gamebase-frontend) 
which provides an easy to use user interface.

## Prerequisites
You will need a [Kubernetes cluster](https://kubernetes.io/) 
(we recommend [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/),
because it is relatively easy to setup) which GameBase will use to manage your game servers.
The GameBase backend is distributed as a [Docker](https://www.docker.com/) image which is built from this repository.
This means you will also need to install Docker on your machine.

Before you continue, please make sure you have both a Kubernetes cluster and Docker installed.

## Deployment
Deployment is quite simple. Just run the docker image for the GameBase backend with the following command:

    docker run -p "$MY_GAMEBASE_BACKEND_PORT:80" -e TZ=Europe/Berlin --mount "type=bind,source=$MY_KUBECONFIG_PATH,target=/root/.kube/config,readonly" --name my-gamebase-backend gamebaseproject/backend

You will need to choose a value for $MY_GAMEBASE_BACKEND_PORT (the port where you want to backend server to run),
$MY_KUBECONFIG_PATH (the path to your kubeconfig file) and maybe you want to change the name of the docker container.
You should see a hash like `b09cf4152900457843a2adb398d330c33bddd87b407036c1d2aa5be3dafe4e05`.
Check if the server is running by running `docker logs my-gamebase-backend`. You should see something like 

    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
    
    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
     - using env:   export GIN_MODE=release
     - using code:  gin.SetMode(gin.ReleaseMode)
    
    [GIN-debug] GET    /                         --> gamebase-daemon/openapi.Index (4 handlers)
    [GIN-debug] POST   /auth/login               --> gamebase-daemon/openapi.AuthLoginPost (4 handlers)
    [GIN-debug] POST   /gs/configure/:id         --> gamebase-daemon/openapi.ConfigureContainer (4 handlers)
    [GIN-debug] DELETE /gs/destroy/:id           --> gamebase-daemon/openapi.DeleteContainer (4 handlers)
    [GIN-debug] POST   /gs/deploy                --> gamebase-daemon/openapi.DeployContainer (4 handlers)
    [GIN-debug] GET    /gs/status                --> gamebase-daemon/openapi.GetStatus (4 handlers)
    [GIN-debug] GET    /gs/templates             --> gamebase-daemon/openapi.ListTemplates (4 handlers)
    [GIN-debug] GET    /gs/restart/:id           --> gamebase-daemon/openapi.RestartContainer (4 handlers)
    [GIN-debug] GET    /gs/start/:id             --> gamebase-daemon/openapi.StartContainer (4 handlers)
    [GIN-debug] GET    /gs/stop/:id              --> gamebase-daemon/openapi.StopContainer (4 handlers)
    [GIN-debug] Listening and serving HTTP on :80

Your output might differ as this is from a debug build.

## Building
You can build this project yourself.
Because the server is written in Go you will need to have [Go](https://golang.org/) installed.
Once setup, run from the repository root:

    make build

This will put the compiled server binary in the directory `out`.

All the other interesting commands are in the Makefile.
