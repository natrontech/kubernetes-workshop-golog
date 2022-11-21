# Scaling Applications
In this tutorial, we are going to show you how to scale applications on Kubernetes. 
Furthermore, we show you how Kubernetes makes sure that the number of requested Pods is up and running and how an application can tell the platform that it is ready to receive requests.

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    export URL=${NAMESPACE}.k8s.golog.ch
    echo $URL
    ```

!!! note

    This tutorial is based on previous tutorials. If you haven’t done so, please complete the following tutorials:
    
    - [Deploying Containers](./deploying-containers.md)
    - [Exposing Services](./exposing-services.md)

## :octicons-tasklist-16: **Task 1**: Scale the test-webserver application
Create a new Deployment in your Namespace. 
So again, lets use the Deployment `deployment.yaml` which we created in the previous tutorial [Deploying Containers](./deploying-containers.md).
Make sure that everything is up and running by using `kubectl get pods --namespace $NAMESPACE`.

If we want to scale our example application, we have to tell the Deployment that we want to have three running replicas instead of one. 
Let’s have a closer look at the existing ReplicaSet:

```bash
kubectl get replicasets --namespace $NAMESPACE
```

Which will give you an output similar to this:

```bash
NAME                        DESIRED   CURRENT   READY   AGE
test-webserver-6564f9788b   1         1         1       54s
```

Or for even more details:

```bash
kubectl get replicaset <replicaset> -o yaml --namespace $NAMESPACE
```

The ReplicaSet shows how many instances of a Pod are desired, current and ready.

Now we scale our application to three replicas:

```bash
kubectl scale deployment test-webserver --replicas=3 --namespace $NAMESPACE
```

Check the number of desired, current and ready replicas:

```bash
kubectl get replicasets --namespace $NAMESPACE
```

```bash
NAME                        DESIRED   CURRENT   READY   AGE
test-webserver-6564f9788b   3         3         3       3m23s
```

Look at how many Pods there are:

```bash
kubectl get pods --namespace $NAMESPACE
```

Which gives you an output similar to this:

```bash
NAME                              READY   STATUS    RESTARTS   AGE
test-webserver-6564f9788b-8np9n   1/1     Running   0          3m54s
test-webserver-6564f9788b-8tzt7   1/1     Running   0          2m2s
test-webserver-6564f9788b-msfvz   1/1     Running   0          2m2s
```

!!! note

    Kubernetes even supports [autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/).

Scaling of Pods is fast as Kubernetes simply creates new containers.

You can check the availability of your Service while you scale the number of replicas up and down in your browser: `https://test.k8s.golog.ch/`.

Now, execute the corresponding loop command for your operating system in another console.

**Linux/MacOS**

```bash
while true; do sleep 1; curl -s https://${URL}/pod/; date "+ TIME: %H:%M:%S,%3N"; done
```

**Windows**

```powershell
while(1) {
  Start-Sleep -s 1
  Invoke-RestMethod https://<namespace>.k8s.golog.ch/pod/
  Get-Date -Uformat "+ TIME: %H:%M:%S,%3N"
}
```

Scale from 3 replicas to 1. The output shows which Pod is still alive and is responding to requests.

??? example "solution"

    ```bash
    kubectl scale deployment test-webserver --replicas=1 --namespace $NAMESPACE
    ```

The requests get distributed amongst the three Pods. As soon as you scale down to one Pod, there should be only one remaining Pod that responds.

Let’s make another test: What happens if you start a new Deployment while our request generator is still running?

```bash
kubectl rollout restart deployment test-webserver --namespace $NAMESPACE
```

During a short period of time, you will see that we won’t get any response from the Service:

```
test-webserver-6564f9788b-8np9n TIME: 15:13:40,249
test-webserver-6564f9788b-8np9n TIME: 15:13:41,499
test-webserver-6564f9788b-8np9n TIME: 15:13:42,719
test-webserver-6564f9788b-8np9n TIME: 15:13:43,945
test-webserver-6564f9788b-8np9n TIME: 15:13:45,190
# no response
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:51,422
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:52,635
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:53,854
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:55,078
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:56,322
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:57,548
test-webserver-5f8bf9b644-5ltlh TIME: 15:13:58,759
```

In our example, we use a very lightweight Pod. 
If we had used a more heavyweight Pod that needed a longer time to respond to requests, we would of course see a larger gap. 
An example for this would be a Java application with a startup time of 30 seconds:

```
test-spring-boot-2-73aln TIME: 16:48:25,251
test-spring-boot-2-73aln TIME: 16:48:26,305
test-spring-boot-2-73aln TIME: 16:48:27,400
test-spring-boot-2-73aln TIME: 16:48:28,463
test-spring-boot-2-73aln TIME: 16:48:29,507
<html><body><h1>503 Service Unavailable</h1>
No server is available to handle this request.
</body></html>
 TIME: 16:48:33,562
<html><body><h1>503 Service Unavailable</h1>
No server is available to handle this request.
</body></html>
 TIME: 16:48:34,601
 ...
test-spring-boot-3-tjdkj TIME: 16:49:20,114
test-spring-boot-3-tjdkj TIME: 16:49:21,181
test-spring-boot-3-tjdkj TIME: 16:49:22,231
```

It is even possible that the Service gets down, and the routing layer responds with the status code 503 as can be seen in the example output above.

In the following chapter we are going to look at how a Service can be configured to be highly available.

### Uninterruptible Deployments
The [rolling update strategy](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/) makes it possible to deploy Pods without interruption. The rolling update strategy means that the new version of an application gets deployed and started. As soon as the application says it is ready, Kubernetes forwards requests to the new instead of the old version of the Pod, and the old Pod gets terminated.

Additionally, [container health checks](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) help Kubernetes to precisely determine what state the application is in.

Basically, there are two different kinds of checks that can be implemented:

Liveness probes are used to find out if an application is still running
Readiness probes tell us if the application is ready to receive requests (which is especially relevant for the above-mentioned rolling updates)
These probes can be implemented as HTTP checks, container execution checks (the execution of a command or script inside a container) or TCP socket checks.

In our example, we want the application to tell Kubernetes that it is ready for requests with an appropriate readiness probe. Our example application has a health check context named health: `https:///<namespace>.k8s.golog.ch/health`

## :octicons-tasklist-16: **Task 2**: Availability during deployment
In our deployment configuration inside the rolling update strategy section, we define that our application always has to be available during an update: `maxUnavailable: 0`

You can directly edit the deployment (or any resource) with:

```bash
kubectl edit deployment test-webserver --namespace $NAMESPACE
```

!!! note

    If you’re not comfortable with `vi` then you can switch to another editor by setting the environment variable `EDITOR` or `KUBE_EDITOR`, e.g. `export KUBE_EDITOR=nano`.

Look for the following section and change the value for maxUnavailable to 0:

```yaml
...
spec:
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 0
    type: RollingUpdate
...
```

Now insert the readiness probe at .spec.template.spec.containers above the resources: {} line:

```yaml
...
spec:
  template:
    spec:
      containers:
      - name: test-webserver
        image: gologch/test-webserver:1.0.0
        # start to copy here
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        # stop to copy here
        resources: {}
...
```

We are now going to verify that a redeployment of the application does not lead to an interruption.

Set up the loop again to periodically check the application’s response (you don’t have to set the $URL variable again if it is still defined):

**Linux/MacOS**

```bash
URL=test.k8s.golog.ch
while true; do sleep 1; curl -s https://${URL}/pod/; date "+ TIME: %H:%M:%S,%3N"; done
```

**Windows**

```powershell
while(1) {
  Start-Sleep -s 1
  Invoke-RestMethod https://test.k8s.golog.ch/pod/
  Get-Date -Uformat "+ TIME: %H:%M:%S,%3N"
}
```

Start a new deployment by editing it (the so-called *ConfigChange* trigger creates the new Deployment automatically):

```bash
kubectl patch deployment test-webserver --patch "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}" --namespace $NAMESPACE
```

### Self-healing
Via the Replicaset we told Kubernetes how many replicas we want. So what happens if we simply delete a Pod?

Look for a running Pod (status `RUNNING`) that you can bear to kill via `kubectl get pods`.

Show all Pods and watch for changes:

```bash
kubectl get pods -w --namespace $NAMESPACE
```

Now delete a Pod (in another terminal) with the following command:

```bash
kubectl delete pod <pod-name> --namespace $NAMESPACE
```

You should see that Kubernetes immediately starts a new Pod to replace the deleted one.
