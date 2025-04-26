# mini-reloader

**Mini Reloader** is a lightweight Kubernetes operator that automatically restarts deployments when their referenced ConfigMaps change.

It ensures your application pods always reflect the latest configuration without manual intervention or GitOps hacks.

## Features
 
- Event-based ConfigMap watching

- Auto restarts Deployments referencing changed ConfigMaps

- Written using Operator SDK and controller-runtime

- No need for custom resources (CRDs)

- Can run with minimal RBAC (only needs access to watch & patch Deployments and ConfigMaps)

## How It Works

1. Watches ConfigMaps and Secrets with label `mini-reloader.pelabs/enabled = true` in the cluster.

2. On change, it:
    - Lists all Deployments in the same namespace.

    - Checks if any of them reference the updated ConfigMap  (via volume or envFrom).

    - If match found, patches the Deployment with a new label like:

        ````yaml
            metadata:
              labels:
                restartedAt: <current-time>
        ````
      
    - This triggers a rolling update automatically.


## Future Improvements

- Support for Secrets also

- Support StatefulSets and DaemonSets

- Add support for projected volumes

- Validate integrity using hashes instead of annotations

- Expose Prometheus metrics (e.g., restarts triggered)

- Optional dry-run mode for GitOps-safe environments

- Multi-namespace or cluster-scoped support

- Send alerts when restarts are triggered



### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

```sh
docker image build . -t <some-registry>/mini-reloader:tag
```

Deploy operator on a cluster
```sh
kubectl apply -f config/mini-reloader-deploy.yaml
kubectl apply -f config/role_binding.yaml
```

Test Operator
````sh
kubectl apply -f config/test-deploy.yaml

````