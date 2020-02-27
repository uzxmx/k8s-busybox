# k8s-busybox

This project provides several tools to help you when using kubernetes cluster. Provided
tools include as below:

* [k8s-eventgenerator](#k8s-eventgenerator)
* [k8s-tlsinfo](#k8s-tlsinfo)

## Getting started

You can download latest prebuilt release [here](https://github.com/uzxmx/k8s-busybox/releases).

## How to build

Generate executables inside `bin/` by:

```sh
$ make
```

Generate a release with cross compilation:

```sh
$ make dist
```

## Tools

### k8s-eventgenerator

**k8s-eventgenerator** is an utility that can help you generate fake events in kubernetes
cluster, especially useful when you use event-based tools like brigade.

### k8s-tlsinfo

**k8s-tlsinfo** is an utility that can help you get the secret tls information in kubernetes
cluster, e.g. certificate common name, expiration.

## License

[MIT License](LICENSE)
