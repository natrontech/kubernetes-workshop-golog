# Sidecar Containers

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

Let’s first have another look at the Pod’s description on the [Kubernetes documentation page](https://kubernetes.io/docs/concepts/workloads/pods/pod/):

!!! quote

    A Pod (as in a pod of whales or pea pod) is a group of one or more containers (such as Docker containers), with shared storage/network, and a specification for how to run the containers. A Pod’s contents are always co-located and co-scheduled, and run in a shared context. A Pod models an application-specific “logical host” - it contains one or more application containers which are relatively tightly coupled — in a pre-container world, being executed on the same physical or virtual machine would mean being executed on the same logical host. The shared context of a Pod is a set of Linux namespaces, cgroups, and potentially other facets of isolation - the same things that isolate a Docker container. Within a Pod’s context, the individual applications may have further sub-isolations applied.

A sidecar container is a utility container in the Pod. Its purpose is to support the main container. It is important to note that the standalone sidecar container does not serve any purpose, it must be paired with one or more main containers. Generally, sidecar containers are reusable and can be paired with numerous types of main containers.

In a sidecar pattern, the functionality of the main container is extended or enhanced by a sidecar container without strong coupling between the two. Although it is always possible to build sidecar container functionality into the main container, there are several benefits with this pattern:

- Different resource profiles, i.e. independent resource accounting and allocation
- Clear separation of concerns at packaging level, i.e. no strong coupling between containers
- Reusability, i.e., sidecar containers can be paired with numerous “main” containers
- Failure containment boundary, making it possible for the overall system to degrade gracefully
- Independent testing, packaging, upgrade, deployment and if necessary rollback

## :octicons-tasklist-16: **Task 1**: Create a sidecar container
We want to create a Pod that runs a sidecar container that reads a log file and prints it to stdout. The main container writes a log file every seconds.

Create a file called `sidecar-pod.yaml` with the following content:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: sidecar-demo
spec:
  initContainers:
    - name: init
      image: busybox
      command:
        - 'sh'
        - '-c'
        - 'echo "init" >> /work-dir/test.log'
      volumeMounts:
        - name: workdir
          mountPath: /work-dir
  containers:
  - name: logwriter
    image: busybox
    command:
      - 'sh'
      - '-c'
      - 'while true; do printf "%s %s\n" "$(date)" >> /work-dir/test.log; sleep 1; done'
    resources:
      requests:
        cpu: 10m
        memory: 16Mi
      limits:
        cpu: 20m
        memory: 32Mi
    volumeMounts:
      - name: workdir
        mountPath: /work-dir
  - name: logreader
    image: busybox
    command:
    - 'sh'
    - '-c'
    - 'while true; do cat /work-dir/test.log; sleep 1; done'
    resources:
      requests:
        cpu: 10m
        memory: 16Mi
      limits:
        cpu: 20m
        memory: 32Mi
    volumeMounts:
    - name: workdir
      mountPath: /work-dir
  volumes:
  - name: workdir
    emptyDir: {}
```

Apply the file:

```bash
kubectl apply -f sidecar-pod.yaml --namespace $NAMESPACE
```

Check the logs of the `logreader` container:

```bash
kubectl logs sidecar-demo logreader --namespace $NAMESPACE --follow
```

You should see the log file being printed to stdout every second.

!!! note

    The `--follow` flag tells `kubectl` to keep the connection open and stream the logs.