# KutEx - Kubernetes to External Service

This little tool helps you to replace a service with one that points to a (local) endpoint
quickly.
It can be helpful especially during local development, when you want services
inside the cluster to access a local running application.

## Usage
To point a service to an external one, simply run the following command:

**Please note that the existing service gets deleted during this process, so make a backup!**

```
kutex <your service name> <your endpoint ip/host>

# E.g.
kutex superservice 192.168.64.1
```