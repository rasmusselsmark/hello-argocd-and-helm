# hello-argocd-and-helm

## Create `kind` cluster

```
brew install kind
kind create cluster
kubectl cluster-info
```


## Create simple Golang app

```
mkdir -p cmd/hello-world-app
touch cmd/hello-world-app/main.go
```

`cmd/hello-world-app/main.go`:
```
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, world!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server is listening on :80...")
    http.ListenAndServe(":80", nil)
}
```

Try it with
```
go build -o hello-world-app cmd/hello-world-app/*
./hello-world-app
```

You should now be able to open http://localhost and see a friendly message from our app.


### Dockerfile

Stop the app if running.

```
mkdir docker
touch docker/Dockerfile
```

`docker/Dockerfile`:
```
FROM golang:1.17-alpine as builder

WORKDIR /app

COPY cmd/ cmd/

RUN go build -o hello-world-app cmd/hello-world-app/*

# Final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/hello-world-app .

EXPOSE 8080

CMD ["./hello-world-app"]
```

Build and verify locally
```
docker build --tag hello-world-app --file docker/Dockerfile .
docker run --publish 80:80 hello-world-app
```

Again, app can be accessed on http://localhost


## Create Helm charts

```
helm create charts
kubectl create namespace hello-world-app
helm install hello-world-app --namespace hello-world-app ./charts 
```

Run the instructions printed from `NOTES.txt`, something like
```
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace hello-world-app -l "app.kubernetes.io/name=charts,app.kubernetes.io/instance=hello-world-app" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace hello-world-app $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace hello-world-app port-forward $POD_NAME 8080:$CONTAINER_PORT
```

Open http://localhost:8080, which is running a standard nginx server


### Deploy our `hello-world-app` using Helm

Make local Docker image available to Kind
```
docker tag hello-world-app kind/hello-world-app:1.0.0
kind load docker-image kind/hello-world-app:1.0.0
```

In `charts/values.yaml` change
```
image:
  repository: kind/hello-world-app:latest
```

And in `charts/Chart.yaml`:
```
appVersion: "1.0.0"
```

Uninstall and deploy again using
```
helm uninstall hello-world-app --namespace hello-world-app
helm install hello-world-app --namespace hello-world-app ./charts
```

Set current `kubectl` namespace and verify state of pod
```
kubectl config set-context --current --namespace=hello-world-app
kubectl get pods
```

If pod is running, we should now be able to access it on http://localhost:8080


### Some technical troubleshooting commands, if pod doesn't start

Kind nodes runs as docker containers
```
$ docker ps
CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS         PORTS                       NAMES
ccba643ce4ff   kindest/node:v1.27.1   "/usr/local/bin/entr…"   2 minutes ago   Up 2 minutes   127.0.0.1:50154->6443/tcp   kind-control-plane
```

Connect to container and list images available to the kind cluster (loaded by `kind load docker-image` command above)
```
docker exec -it kind-control-plane bash
root@kind-control-plane:/# crictl images
```

Get more information about pod state
```
kubectl get pods
kubectl describe pod <name>
```


### Useful `helm` commands

```
helm ls --namespace hello-world-app
```


### Useful `kubectl` commands

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

```
kubectl get namespaces
kubectl get pods --all-namespaces
```


## Install and run ArgoCD

```
kubectl create namespace argocd
kubectl apply --namespace argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
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


### Use port forwarding to open ArgoCD
```
kubectl get services -n argocd
kubectl port-forward service/argocd-server -n argocd 8080:443
```

ArgoCD can now be opened at https://localhost:8080 (you need to accept the self-signed certificate)

Get password for user `admin`:
```
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```


## Create ArgoCD app (manually using UI)

1. Click `+ New app` button
1. Enter
  - Application name: hello-world-app
  - Project name: default
  - Sync policy: Manual
  - Auto-create namespace: On
  - Repository url: https://github.com/rasmusselsmark/hello-argocd-and-helm
  - Path: charts
  - Cluster URL: https://kubernetes.default.svc
  - Namespace: hello-world-app-argocd
  - Values: values.yaml
1. You should now be able to see the app. Click Sync to synchronize it

### Use port forwarding to access our app running in `hello-world-app-argocd` namespace

```
kubectl get services -n hello-world-app-argocd
kubectl port-forward service/hello-world-app-charts -n hello-world-app-argocd 8081:80
```

You can now access the app at http://localhost:8081/


## Create ArgoCD app using `kubectl`

Instead of creating manually, we can apply the manifest using
```
kubectl apply -n argocd -f apps/hello-world-app.yaml
```

Note this manifest has auto-sync enabled, so no need to manually refresh in ArgoCD UI.


## Using ArgoCD App of Apps Pattern

The [ArgoCD App of Apps pattern](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/#app-of-apps-pattern) allows us to configure ArgoCD as code, which is particularly useful when having more than one application.

Video: [The App of Apps Pattern for Argo CD in Less than 8 Minutes](https://www.youtube.com/watch?v=2pvGL0zqf9o)

By running following command, we're configuring ArgoCD to automatically look in the `apps/` folder for applications:
```
kubectl apply -n argocd -f app-of-apps/app-of-apps.yaml
```
