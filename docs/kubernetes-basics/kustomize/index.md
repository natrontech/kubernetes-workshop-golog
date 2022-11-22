# Kustomize

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

[Kustomize](https://kustomize.io/) is a tool to manage YAML configurations for Kubernetes objects in a declarative and reusable manner. In this lab, we will use Kustomize to deploy the same app for two different environments.

## Installation
Kustomize can be used in two different ways:

- As a standalone `kustomize` binary, downloadable from [here](https://kubernetes-sigs.github.io/kustomize/installation/)
- With the parameter `--kustomize` or `-k` in certain `kubectl` subcommands such as `apply` or `create`

!!! note

    You might get a different behaviour depending on which variant you use. The reason for this is that the version built into `kubectl` is usually older than the standalone binary.

## Usage
The main purpose of Kustomize is to build configurations from a predefined file structure (which will be introduced in the next section):

```bash
kustomize build <path-to-kustomization-directory>
```

The same can be achieved with kubectl:

```bash
kubectl apply -k <path-to-kustomization-directory>
```

The next step is to apply this configuration to the Kubernetes cluster:

```bash
kustomize build <path-to-kustomization-directory> | kubectl apply -f -
```

Or in one kubectl command with the parameter `-k` instead of `-f`:

```bash
kubectl apply -k <path-to-kustomization-directory>
```

## :octicons-tasklist-16: **Task 1**: Prepare a Kustomize config
We are going to deploy a simple application:

- The Deployment starts an application based on nginx
- A Service exposes the Deployment
- The application will be deployed for two different example environments, integration and production

Kustomize allows inheriting Kubernetes configurations. We are going to use this to create a base configuration and then override it for the different environments. Note that Kustomize does not use templating. Instead, smart patch and extension mechanisms are used on plain YAML manifests to keep things as simple as possible.

### File structure
The structure of a Kustomize configuration typically looks like this:

```
.
├── base
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays
    ├── production
    │   ├── deployment-patch.yaml
    │   ├── kustomization.yaml
    │   └── service-patch.yaml
    └── staging
        ├── deployment-patch.yaml
        ├── kustomization.yaml
        └── service-patch.yaml
```

### Base
Let’s have a look at the `base` directory first which contains the base configuration. There’s a `deployment.yaml` with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kustomize-app
spec:
  selector:
    matchLabels:
      app: kustomize-app
  template:
    metadata:
      labels:
        app: kustomize-app
    spec:
      containers:
        - name: kustomize-app
          image: ghcr.io/natrongmbh/kubernetes-workshop-golog-test-webserver:latest
          env:
            - name: APPLICATION_NAME
              value: app-base
          command:
            - sh
            - -c
            - |-
              set -e
              /bin/echo "My name is $APPLICATION_NAME"
              /usr/local/bin/go              
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
```

There’s also a Service for our Deployment in the corresponding `base/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: kustomize-app
spec:
  ports:
    - port: 80
      targetPort: 80
  selector:
    app: kustomize-app
```

And there’s an additional `base/kustomization.yaml` which is used to configure Kustomize:

```yaml
resources:
  - service.yaml
  - deployment.yaml
```

It references the previous manifests `service.yaml` and `deployment.yaml` and makes them part of our base configuration.

### Overlays
Now let’s have a look at the other directory which is called `overlays`. It contains two subdirectories `staging` and `production` which both contain a `kustomization.yaml` with almost the same content.

`overlays/staging/kustomization.yaml`:

```yaml
nameSuffix: -staging
bases:
  - ../../base
patchesStrategicMerge:
  - deployment-patch.yaml
  - service-patch.yaml
```

`overlays/production/kustomization.yaml`:

```yaml
nameSuffix: -production
bases:
  - ../../base
patchesStrategicMerge:
  - deployment-patch.yaml
  - service-patch.yaml
```

Only the first key `nameSuffix` differs.

In both cases, the `kustomization.yaml` references our base configuration. However, the two directories contain two different `deployment-patch.yaml` files which patch the `deployment.yaml` from our base configuration.

`overlays/staging/deployment-patch.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kustomize-app
spec:
  selector:
    matchLabels:
      app: kustomize-app-staging
  template:
    metadata:
      labels:
        app: kustomize-app-staging
    spec:
      containers:
        - name: kustomize-app
          env:
            - name: APPLICATION_NAME
              value: kustomize-app-staging
```

`overlays/production/deployment-patch.yaml`:
    
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kustomize-app
spec:
  selector:
    matchLabels:
      app: kustomize-app-production
  template:
    metadata:
      labels:
        app: kustomize-app-production
    spec:
      containers:
        - name: kustomize-app
          env:
            - name: APPLICATION_NAME
              value: kustomize-app-production
```

The main difference here is that the environment variable `APPLICATION_NAME` is set differently. The `app` label also differs because we are going to deploy both Deployments into the same Namespace.

The same applies to our Service. It also comes in two customizations so that it matches the corresponding Deployment in the same Namespace.

`overlays/staging/service-patch.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: kustomize-app
spec:
  selector:
    app: kustomize-app-staging
```

`overlays/production/service-patch.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: kustomize-app
spec:
  selector:
    app: kustomize-app-production
```

!!! info

    All files mentioned above are also directly accessible from GitHub.

Prepare the files as described above in a local directory of your choice.

## :octicons-tasklist-16: **Task 2**: Deploy with Kustomize
We are now ready to deploy both apps for the two different environments. For simplicity, we will use the same Namespace.

```bash
kubectl apply -k overlays/staging --namespace $NAMESPACE
```

```bash
kubectl apply -k overlays/production --namespace $NAMESPACE
```

As you can see, we now have two deployments and services deployed. Both of them use the same base configuration. However, they have a specific configuration on their own as well.

Let’s verify this. Our app writes a corresponding log entry that we can use for analysis:

```bash
kubectl get pods --namespace $NAMESPACE
```

```bash
kubectl logs <pod-name> --namespace $NAMESPACE
```

## Further Reading
Kustomize has more features of which we just covered a couple. Please refer to the docs for more information.

- Kustomize documentation: [https://kubernetes-sigs.github.io/kustomize/](https://kubernetes-sigs.github.io/kustomize/)
- API reference: [https://kubernetes-sigs.github.io/kustomize/api-reference/](https://kubernetes-sigs.github.io/kustomize/api-reference/)
- Another kustomization.yaml reference: [https://kubectl.docs.kubernetes.io/pages/reference/kustomize.html](https://kubectl.docs.kubernetes.io/pages/reference/kustomize.html)
- Examples: [https://github.com/kubernetes-sigs/kustomize/tree/master/examples](https://github.com/kubernetes-sigs/kustomize/tree/master/examples)
