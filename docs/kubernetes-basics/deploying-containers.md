# Deploying Containers
In this tutorial, we are going to deploy our first container image and look at the concepts of Pods, Services, and Deployments.

## :octicons-tasklist-16: **Task 1**: Start and stop a single Pod
After we’ve familiarized ourselves with the platform, we are going to have a look at deploying a pre-built container image or any other public container registry.

First, we are going to directly start a new Pod.
For this we have to define our Kubernetes Pod resource definition. 
Create a new file `03_pod.yaml` with the following content:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-webserver
spec:
  containers:
  - image: ghcr.io/natrongmbh/kubernetes-workshop-test-webserver:latest
    imagePullPolicy: Always
    name: test-webserver
    resources:
      limits:
        cpu: 20m
        memory: 32Mi
      requests:
        cpu: 10m
        memory: 16Mi
```

Now we can apply this with:

```bash
kubectl apply -f 03_pod.yaml --namespace <namespace>
```

The output should be:
```bash
pod/test-webserver created
```

Use `kubectl get pods --namespace <namespace>` in order to show the running Pod:

```bash
kubectl get pods --namespace <namespace>
```

Which gives you an output similar to this:

```bash
NAME                READY   STATUS    RESTARTS   AGE
test-webserver   1/1     Running   0          2m
```

Now we can delete the Pod with:

```bash
kubectl delete pod test-webserver --namespace <namespace>
```

## :octicons-tasklist-16: **Task 2**: Create a Deployment
In some use cases it can make sense to start a single Pod. But this has its downsides and is not really a common practice. Let’s look at another concept which is tightly coupled with the Pod: the so-called *Deployment*. A Deployment ensures that a Pod is monitored and checks that the number of running Pods corresponds to the number of requested Pods.

To create a new Deployment we first define our Deployment in a new file `03_deployment.yaml` with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: test-webserver
  name: test-webserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-webserver
  template:
    metadata:
      labels:
        app: test-webserver
    spec:
      containers:
      - image: ghcr.io/natrongmbh/kubernetes-workshop-test-webserver:latest
        name: test-webserver
        resources:
          requests:
            cpu: 10m
            memory: 16Mi
          limits:
            cpu: 20m
            memory: 32Mi
```

And with this we create our Deployment inside our already created namespace:

```bash
kubectl apply -f 03_deployment.yaml --namespace <namespace>
```

The output should be:

```bash
deployment.apps/test-webserver created
```

Kubernetes creates the defined and necessary resources, pulls the container image (in this case from ghcr.io) and deploys the Pod.

Use the command `kubectl get` with the `-w` parameter in order to get the requested resources and afterward watch for changes.

!!! note

    The `kubectl get -w `command will never end unless you terminate it with `CTRL-c`.

```bash
kubectl get pods -w --namespace <namespace>
```

!!! note
    
    Instead of using the `-w` parameter you can also use the `watch` command which should be available on most Linux distributions:

    ```bash
    watch kubectl get pods --namespace <namespace>
    ```

This process can last for some time depending on your internet connection and if the image is already available locally.

!!! note

    If you want to create your own container images and use them with Kubernetes, you definitely should have a look at [these best practices](https://docs.openshift.com/container-platform/4.11/openshift_images/create-images.html) and apply them. This image creation guide may be for OpenShift, however it also applies to Kubernetes and other container platforms.

### Creating Kubernetes resources
There are two fundamentally different ways to create Kubernetes resources. You’ve already seen one way: Writing the resource’s definition in YAML (or JSON) and then applying it on the cluster using `kubectl apply`.

The other variant is to use helper commands. These are more straightforward: You don’t have to copy a YAML definition from somewhere else and then adapt it. However, the result is the same. The helper commands just simplify the process of creating the YAML definitions.

As an example, let’s look at creating above deployment, this time using a helper command instead. If you already created the Deployment using above YAML definition, you don’t have to execute this command:

```bash
kubectl create deployment test-webserver --image=ghcr.io/natrongmbh/kubernetes-workshop-test-webserver:latest --namespace <namespace>
```

It’s important to know that these helper commands exist. However, in a world where GitOps concepts have an ever-increasing presence, the idea is not to constantly create these resources with helper commands. Instead, we save the resources’ YAML definitions in a Git repository and leave the creation and management of those resources to a tool.

## :octicons-tasklist-16: **Task 3**: Viewing the created resources
Display the created Deployment using the following command:

```bash
kubectl get deployments --namespace <namespace>
```

A [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) defines the following facts:

- Update strategy: How application updates should be executed and how the Pods are exchanged
- Containers
    - Which image should be deployed
    - Environment configuration for Pods
    - ImagePullPolicy
- The number of Pods/Replicas that should be deployed

By using the `-o` (or `--output`) parameter we get a lot more information about the deployment itself. You can choose between YAML and JSON formatting by indicating `-o yaml` or `-o json`. In this training we are going to use YAML, but please feel free to replace `yaml` with `json` if you prefer.

```bash
kubectl get deployments test-webserver -o yaml --namespace <namespace>
```

After the image has been pulled, Kubernetes deploys a Pod according to the Deployment:

```bash
kubectl get pods --namespace <namespace>
```

which gives you an output similar to this:

```bash
NAME                                READY   STATUS    RESTARTS   AGE
test-webserver-7f7b9b9b4-2j2xg      1/1     Running   0          2m
```

The Deployment defines that one replica should be deployed — which is running as we can see in the output. This Pod is not yet reachable from outside the cluster.
