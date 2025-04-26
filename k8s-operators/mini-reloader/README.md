# mini-reloader

**Mini Reloader** is a lightweight Kubernetes operator that automatically restarts deployments when their referenced ConfigMaps or Secrets change.

It ensures your application pods always reflect the latest configuration without manual intervention or GitOps hacks.

## Features
 
- Event-based ConfigMap/Secret watching

- Auto restarts Deployments referencing changed ConfigMaps or Secrets

- Written using Operator SDK and controller-runtime

- No need for custom resources (CRDs)

- Can run with minimal RBAC (only needs access to watch & patch Deployments and ConfigMaps)

## How It Works

1. Watches ConfigMaps and Secrets with label `mini-reloader.pelabs/enabled = true` in the cluster.

2. On change, it:
    - Lists all Deployments in the same namespace.

    - Checks if any of them reference the updated ConfigMap or Secret (via volume or envFrom).

    - If match found, patches the Deployment with a new label like:

        ````yaml
            metadata:
              labels:
                restartedAt: <current-time>
        ````
      
    - This triggers a rolling update automatically.


## Future Improvements

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
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/mini-reloader:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/mini-reloader:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/mini-reloader:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/mini-reloader/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

