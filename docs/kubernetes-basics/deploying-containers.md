# Deploying Containers
In this tutorial, we are going to deploy our first container image and look at the concepts of Pods, Services, and Deployments.

## :octicons-tasklist-16: **Task**: Start and stop a single Pod
After we’ve familiarized ourselves with the platform, we are going to have a look at deploying a pre-built container image or any other public container registry.

First, we are going to directly start a new Pod.
For this we have to define our Kubernetes Pod resource definition. 
Create a new file `03_pod.yaml` with the following content:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-webserver-v1
spec:
  containers:
  - image: ghcr.io/natrongmbh/kubernetes-workshop-test-webserver-v1:latest
    imagePullPolicy: Always
    name: test-webserver-v1
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
pod/test-webserver-v1 created
```

Use `kubectl get pods --namespace <namespace>` in order to show the running Pod:

```bash
kubectl get pods --namespace <namespace>
```

Which gives you an output similar to this:

```bash
NAME                READY   STATUS    RESTARTS   AGE
test-webserver-v1   1/1     Running   0          2m
```

Now we can delete the Pod with:

```bash
kubectl delete pod test-webserver-v1 --namespace <namespace>
```

