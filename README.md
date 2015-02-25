deployer
========

[![Build Status](https://img.shields.io/circleci/project/jsdir/deployer.svg)](https://circleci.com/gh/jsdir/deployer)

## What is it?

Deployer manages and deploys releases for distributed applications. Although it currently only works with kubernetes, it can be easily extended to work with other cloud orchestration systems. Contributions are more than welcome!

Deployer can:

- Manage releases through a CLI
- Deploy to multiple environments and individually apply releases to them
- Use pluggable environment backends

## Quickstart

We'll use docker for demonstrating how deployer works with a kubernetes cluster. First, pull some images.

```bash
$ docker pull jsdir/deployer
$ docker pull llamashoes/dind-kubernetes
```

Start a local kubernetes cluster.

```bash
$ docker run -d -p 127.0.0.1:8888:8888 -p 8091:8091 -p 8092:8092 --privileged llamashoes/dind-kubernetes
```

Next, we'll write some configuration files for deployer.

### Config

The config decalres the environments that we want to deploy to. Here, we'll set up staging and production environments.

```bash
mkdir -p /tmp/deployer-demo
cat <<EOF > /tmp/deployer-demo/config.json
{
  "port": 7654,
  "environments": {
    "production": {
      "type": "kubernetes",
      "manifestGlob": "/tmp/deployer-demo/manifest.json",
      "cmd": "docker kubectl --server=http://localhost:7654"
    }
  }
}
EOF
```

### Kubernetes manifest

This manifest injects information about the deployment into the container's environment through templates. This is done using [deployer-kubernetes](https://github.com/jsdir/deployer-kubernetes).

```bash
cat <<EOF > /tmp/deployer-demo/manifest.json
{
  "id": "server",
  "kind": "Pod",
  "apiVersion": "v1beta1",
  "desiredState": {
    "manifest": {
      "version": "v1beta1",
      "id": "server",
      "containers": [{
        "name": "web-demo",
        "image": "{{.Services.web-demo}}",
        "cpu": 100,
        "ports": [{
          "name": "http",
          "containerPort": 8091,
          "hostPort": 8091
        }],
        "env": [{
          "name": "PORT",
          "value": "8091"
        }]
      }]
    }
  }
}
EOF
```

### Creating releases

Now that the config is set up, start the daemon.

```shell
$ docker run -d -p 7654:7654 --net=host -v /tmp/deployer-demo:/data jsdir/deployer deployerd --config /data/config.json
```

Commands can be sent to `deployerd` through the HTTP API, or through `deployer`, a simple command line interface that's bundled with `jsdir/deployer`.

We'll use the API to create a build. This step would be run manually or in CI once a docker image is tested, built, and uploaded to the registry. In this case, we'll continue as if you just uploaded `jsdir/deployer-web-demo#version1` to the registry.

```shell
$ curl --data "service=web-demo&tag=jsdir/deployer-web-demo#version1" localhost:7654/builds
```

Next, we'll create a release with the new build using the CLI.

```shell
$ alias deployer='docker run --net=host jsdir/deployer deployer --addr="http://localhost:7654"'
$ deployer release web-demo jsdir/deployer-web-demo#version1
{"id": 0, "name": "super-panda", "services": {"web-demo": "jsdir/deployer-web-demo#version1"}}
0
```

This creates a release. The release is a named, atomic mapping of services and their builds. Creating a release will always block until the required service builds exist. Since the
`web-demo` build `jsdir/deployer-web-demo#version1` already exists, the client responds immediately.

### Deploying releases

To deploy this release, we'll use the CLI.

```shell
$ deployer deploy 0 staging
```

This deploys release `0` to the `staging` environment. To verify that it works:

```shell
$ curl http://localhost:8091
Hello from `jsdir/deployer-web-demo#1`! (super-panda)
I'm running in the `staging` environment.
```

Since it works, let's deploy the release to production.

```shell
$ deployer deploy staging production
```

This syntax deploys the current release at the `staging` environment to the `production` environment.

```shell
$ curl http://localhost:8092
Hello from `jsdir/deployer-web-demo#1`! (super-panda)
I'm running in the `production` environment.
```

Replacing existing deployments with new releases is simple:

```shell
$ curl http://localhost:7654/builds service=web-demo&tag=jsdir/deployer-web-demo#2

# Creating a new release extends the existing services.
$ deployer release web-demo jsdir/deployer-web-demo#2
{"id": 2, "name": "bubbly-whale", "services": {"web-demo": "jsdir/deployer-web-demo#2"}}
2

$ deployer deploy 2 staging
$ curl http://localhost:8091
Hello from `jsdir/deployer-web-demo#2`! (bubbly-whale)
I'm running in the `staging` environment.

$ deployer deploy staging production
$ curl http://localhost:8092
Hello from `jsdir/deployer-web-demo#1`! (super-panda)
I'm running in the `production` environment.
```

## Examples and tutorials

New examples are coming soon:

- assets: Service builds are not limited to docker images. This example will show how to use deployer to organized builds for both a web server, and client javascript assets that are on a CDN
- ci: This example will show how to use deployer to set up a continuous deployment workflow with CI servers like [Drone](https://github.com/drone/drone) and [CircleCI](https://circleci.com).

## Documentation

- Concepts and design
- Usage
- API

## Roadmap

- Web frontend
- More pluggable backends
- Rollbacks
- Availability zones
- Integration with Hubot, Slack, and Flowdock
