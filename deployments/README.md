
# Deployment

## Kubernetes

Example Kubernetes manifests are located in the [k8s](./k8s/) subdirectory. These are only useful for small deployments. For larger deployments, use a proper Redis cluster like one deployed with [this Redis operator](https://ot-container-kit.github.io/redis-operator/). 

```
kubectl apply -r -f ./k8s/
```

## Docker Compose

The example [docker-compose.yml](./docker-compose.yml) file can be brought up with `docker-compose`.

```
docker-compose up
```

<!-- vim: set conceallevel=2 et ts=2 sw=2: -->
