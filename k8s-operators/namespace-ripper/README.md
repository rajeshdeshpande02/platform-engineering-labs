# namespace-ripper

**A Kubernetes Operator for Time-To-Live (TTL) based Namespace Cleanup.**
 
 This operator, manages the lifecycle of namespaces based on a
 time-to-live (TTL) value specified in the NamespaceTTL custom resource.

 It performs the following steps:
 1. Fetches the NamespaceTTL instance specified in the request.
 ```
 apiVersion: core.pelabs.com/v1
kind: NamespaceTTL
metadata:
  labels:
    app.kubernetes.io/name: namespace-ripper
    app.kubernetes.io/managed-by: kustomize
  name: namespacettl-sample
spec:
  ttl: 30s
  exceptions:
    - "test"
    - "test2"

 ```
 2. Parses the TTL value from the NamespaceTTL spec.
 3. Lists all namespaces in the cluster.
 4. Deletes namespaces that have existed longer than the specified TTL,
    excluding protected namespaces and any additional exceptions specified
    in the NamespaceTTL spec.
 5. Requeues the reconciliation after 15 seconds to ensure periodic cleanup.


## Benefits of 'namespace-ripper' Operator

🧹 1. Automated Cleanup of Stale Namespaces
- Avoids accumulation of old namespaces created for testing, staging, or CI/CD pipelines.

- Saves cluster resources and improves overall hygiene.

- Reduces manual cleanup tasks for DevOps/Platform teams.

🔐 2. Controlled Exception Mechanism
- You’ve smartly added an exceptions field in the CRD.

- Ensures critical namespaces (test, dev, platform, etc.) are never deleted accidentally.

- This supports Zero Trust principles by reducing the blast radius of misconfigurations.

⏱️ 3. TTL-Based Lifecycle Management
- Enforces namespace lifespan policies (like ttl=30s, ttl=1h, etc.).

- Useful in sandbox environments or preview deployments where namespaces should auto-expire.

- Can be integrated into GitOps workflows, where TTL is defined declaratively.

🧩 4. Custom Resource Driven
- Everything is managed via a single, declarative NamespaceTTL CR.

- Easy for platform teams to onboard this into existing Kustomize/Helm/ArgoCD setups.

- Simple to audit or replicate across clusters.

🧠 5. Fully Kubernetes-Native & Golang-Based
- No external dependencies — runs as a Kubernetes operator.

- Native reconciliation loop ensures consistency and reliability.

- Easy to extend for more features (like dry-run, label selectors, or notifications).

## 🏢 Business Value of `namespace-ripper` Operator

| 💼 Area                  | 🌟 Value                                                                 |
|--------------------------|---------------------------------------------------------------------------|
| **Developer Experience** | Developers can create test namespaces without worrying about cleanup.     |
| **Platform Engineering** | Enforces automated TTL policies, reducing manual toil and human error.    |
| **Security & Compliance**| Prevents namespace sprawl and reduces risk of shadow resources.           |
| **Cost Optimization**    | Deletes unused namespaces and workloads, saving cluster resources.        |
| **SRE/Operations**       | Improves cluster hygiene and supports lifecycle management practices.     |

## Roadmap

Add dry-run mode (spec.dryRun: true) to preview deletions.

Add label/annotation-based filtering to target specific namespaces.

Send Slack or webhook notifications before deletion.

Track deletions in a CR status field (status.deletedNamespaces).

Expose Prometheus metrics for number of namespaces deleted or skipped.

Webhook for CR validation to enforce minimum TTL.

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/namespace-ripper:tag
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
make deploy IMG=<some-registry>/namespace-ripper:tag
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
make build-installer IMG=<some-registry>/namespace-ripper:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https:raw.githubusercontent.com/<org>/namespace-ripper/<tag or branch>/dist/install.yaml
```



