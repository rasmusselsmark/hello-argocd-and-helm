# hello-argocd-and-helm

### Create `kind` cluster

```
brew install kind
kind create cluster --name hello-argo-and-helm
kubectl cluster-info
```

### Install and run ArgoCD

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

See what's running
```
kubectl get all --namespace argocd
```

### Kubernetes types

* Pods: Instances of the application
* Services: Provides networking, e.g. this is where you can see the IP address of a service
* Deployments: Manages replica sets for the application
* ReplicaSets: Ensures a specified number of replicas are running
* StatefulSets: Manages stateful applications by providing stable network identities, ordered deployment, and persistent storage for Pods with unique identities.


### Some `kubectl` commands

See configured clusters/contexts:
```
less ~/.kube/config
```

Get current configured context for `kubectl`:
```
kubectl config current-context
```

Set current context and namespace:
```
kubectl config use-context kind-hello-argo-and-helm
kubectl config set-context --current --namespace=argocd
```


### Forward ports to access ArgoCD
```
kubectl get services -n argocd
kubectl port-forward service/argocd-server -n argocd 8080:443
```

ArgoCD can now be opened at https://localhost:8080 (you need to accept the self-signed certificate)

Get password for user `admin`:
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```


### Create ArgoCD app

1. Click `+ New app` button
1. Enter
  - Application name: hello-argocd-and-helm-dev
  - Project name: default
  - Sync policy: Manual
  - Auto-create namespace: On
  - Repository url: https://github.com/rasmusselsmark/hello-argocd-and-helm
  - 
