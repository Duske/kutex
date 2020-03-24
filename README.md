# KutEx - Kubernetes to External Service
[![Build Status](https://travis-ci.org/Duske/kutex.svg?branch=master)](https://travis-ci.org/Duske/kutex)

This little tool helps you to replace a service with one that points to a (local) endpoint
quickly.
It can be helpful especially during local development, when you want services
inside the cluster to access a local running application.

## Usage
A command-line interface to point services to external applications.

```
Usage:
  kutex [flags]
  kutex [command]

Available Commands:
  help        Help about any command
  replace     Point a service to an external service 
  restore     Restore previously replaced service

Flags:
  -h, --help   help for kutex
```

### Replace
Replace a service with a service pointing to an external endpoint. While the original service gets deleted, it is attached as an annotation to the new service and can be restored later.

```
Usage:
  kutex replace <servicename> <external host> [flags]

Flags:
  -h, --help                help for replace
  -k, --kubeconfig string   (optional) absolute path to the kubeconfig file (default "/Users/johndoe/.kube/config")
  -n, --namespace string    The namespace the current service is placed (default "default")

```

**Example**

```
kutex replace superservice 192.168.64.1
```

### Restore
Restore a replaced service by kutex. This bring the original service back in place.

Alternatively,  *you can also apply your original service resource definition*.


```
Usage:
  kutex restore [flags]

Flags:
  -h, --help                help for restore
  -k, --kubeconfig string   (optional) absolute path to the kubeconfig file (default "/Users/dchabrowski/.kube/config")
  -n, --namespace string    The namespace the current service is placed (default "default")
```

Example

```
kutex restore
```
