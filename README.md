# proxy-go

This is a simple service which takes requests and forwards
them inside a K8s cluster (that's the initial goal). The
motivation is to avoid doing port forward on every pod and
reconnecting each time a pod dies. Relying on the services
in front of those pods and, in the end, port forwarding
just one pod.

## Usage

```shell
kubectl port-forward -n <namespace> <proxy-go-pod> 8080:8080
```

Making a request to a service named `echo` at `/` is as
follows:

```shell
curl http://localhost/echo
```

If the service is serving on a different port from `80`, just
add the port like this:

```shell
curl http://localhost/echo:<port>
```

## K8s

### Basic deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: proxy
spec:
  selector:
    matchLabels:
      app: proxy
  replicas: 1
  template:
    metadata:
      labels:
        app: proxy
    spec:
      containers:
      - image: dlilue/proxy-go
        imagePullPolicy: Always
        name: proxy
        ports:
        - containerPort: 8080

```