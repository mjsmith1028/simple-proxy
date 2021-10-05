# Simple-proxy

## Purpose
Simple-proxy is a simple deployment that proxies requests to K8s. The specific use case that Simple-proxy was designed
for is to expose link-local IP addresses to enable [AAD Pod Identity](https://github.com/Azure/aad-pod-identity) testing
locally with [Telepresence](https://www.telepresence.io/).

K8s supports services without selectors, however, this does not allow link-local IP addresses. See
https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors.

AAD Pod Identity runs a daemonset for the NMI service responsible for requesting tokens on the address 169.254.169.254.
Simple-proxy allows you to proxy requests to 169.254.169.254 when you are port-forwarding, ultimately allowing you
to request tokens from your local machine.

Although Simple-proxy was written for this specific case, it is hopefully generic enough to be used for other 
similar use cases. Feel free to open PRs or Issues for new features.

## Install
- `helm repo add simple-proxy https://mjsmith1028.github.io/simple-proxy/`
- `helm repo update`
- `helm install simple-proxy simple-proxy/simple-proxy`

## Uninstall
- `helm uninstall simple-proxy`

## Customize
The following command will deploy simple-proxy in your desired namespace and set the desired azure identity binding to 
the pod. See helm chart for additional configuration values.
- `helm install simple-proxy simple-proxy/simple-proxy --namespace <your-namespace> --set namespace=<your-namespace>,podLabels.aadpodidbinding=<your-identity>`

## Example
The purpose of this example is to provide a high level guide for exposing AAD Pod Identity in K8s locally.
This assumes you have knowledge of both Telepresence and AAD Pod Identity, and has only been tested on OSX.
- Install Simple-proxy to K8s in namespace demo-ns and the identity demo-user
  - `helm install simple-proxy simple-proxy/simple-proxy --namespace demo-ns --set namespace=demo-ns,podLabels.aadpodidbinding=demo-user`
- Create an alias for 169.254.169.254 to send requests to localhost on your host machine
  - `sudo ifconfig lo0 169.254.169.254 alias`
- Port forward local requests to 169.254.169.254 to Simple-proxy
  - `sudo kubectl port-forward -n demo-ns deployment/simple-proxy 80:8080 --address=169.254.169.254`
- Send a request locally and confirm a token is successfully returned
  - `curl http://169.254.169.254/metadata/identity/oauth2/token/?resource=https://management.core.windows.net/`
- Execeute the telepresence intercept command and debug your application which uses AAD Pod Identity
  - `telepresence intercept <your-app>`
 