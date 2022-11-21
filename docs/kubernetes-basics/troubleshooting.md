# Troubleshooting

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

This tutorial will help you troubleshoot your application and show you some tools that can make troubleshooting easier.

## Logging into a container
Running containers should be treated as immutable infrastructure and should therefore not be modified. 
However, there are some use cases in which you have to log into your running container. 
Debugging and analyzing is one example for this.

## :octicons-tasklist-16: **Task 1**: Shell into Pod
With Kubernetes you can open a remote shell into a Pod without installing SSH by using the command `kubectl exec`. The command can also be used to execute any command in a Pod. With the parameter `-it` you can leave an open connection.

!!! note

    On Windows, you can use Git Bash and `winpty`.

Choose a Pod with `kubectl get pods --namespace $NAMESPACE` and execute the following command:

```bash
kubectl exec -it <pod> --namespace $NAMESPACE -- /bin/sh
```

!!! note

    If Bash is  available in the Pod you can fallback to – `/bin/bash` instead of – `/bin/sh`.

You now have a running shell session inside the container in which you can execute every binary available, e.g.:

```bash
~@<pod>:/# ls -la /usr/local/bin/
total 6308
drwxr-xr-x    1 root     root            16 Nov 19 13:43 .
drwxr-xr-x    1 root     root            17 Aug  9 08:47 ..
-rwxr-xr-x    1 root     root       6456761 Nov 19 13:43 go
```

With `exit` or `CTRL+d` you can leave the container and close the connection:

```bash
~@<pod>:/# exit
```

## :octicons-tasklist-16: **Task 2**: Single commands
Single commands inside a container can also be executed with `kubectl exec`:

```bash
kubectl exec <pod> --namespace $NAMESPACE -- env
```

### Watching log files
Log files of a Pod can be shown with the following command:

```bash
kubectl logs <pod> --namespace $NAMESPACE
```

The parameter `-f` allows you to follow the log file (same as `tail -f`). With this, log files are streamed and new entries are shown immediately.

When a Pod is in state `CrashLoopBackOff` it means that although multiple attempts have been made, no container inside the Pod could be started successfully. Now even though no container might be running at the moment the `kubectl log`s command is executed, there is a way to view the logs the application might have generated. This is achieved using the `-p` or `--previous` parameter:

```bash
kubectl logs <pod> --namespace $NAMESPACE -p
```

## :octicons-tasklist-16: **Task 3**: Port forwarding
Kubernetes allows you to forward arbitrary ports to your development workstation. This allows you to access admin consoles, databases, etc., even when they are not exposed externally. Port forwarding is handled by the Kubernetes control plane nodes and therefore tunneled from the client via HTTPS. This allows you to access the Kubernetes platform even when there are restrictive firewalls or proxies between your workstation and Kubernetes.

Get the name of the Pod:

```bash
kubectl get pods --namespace $NAMESPACE
```

Then execute the port forwarding command using the Pod’s name:

```bash
kubectl port-forward <pod> 8080:8080 --namespace $NAMESPACE
```

!!! note

    Use the additional parameter `--address <IP address>` (where `<IP address>` refers to a NIC’s IP address from your local workstation) if you want to access the forwarded port from outside your own local workstation.

The output of the command should look like this:

```
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

Don’t forget to change the Pod name to your own installation. If configured, you can use auto-completion.

The application is now available with the following link: `http://localhost:8080/` . Or try a curl command:

```bash
curl http://localhost:8080/
```

With the same concept you can access databases from your local workstation or connect your local development environment via remote debugging to your application in the Pod.

[This documentation page](https://kubernetes.io/docs/tasks/access-application-cluster/port-forward-access-application-cluster) offers some more details about port forwarding.

!!! note

    The `kubectl port-forward` process runs as long as it is not terminated by the user. So when done, stop it with `CTRL-c`.

## Events
Kubernetes maintains an event log with high-level information on what’s going on in the cluster. It’s possible that everything looks okay at first but somehow something seems stuck. Make sure to have a look at the events because they can give you more information if something is not working as expected.

Use the following command to list the events in chronological order:

```bash
kubectl get events --sort-by=.metadata.creationTimestamp --namespace $NAMESPACE
```

## Dry-run
To help verify changes, you can use the optional `kubectl` flag `--dry-run=client -o yaml` to see the rendered YAML definition of your Kubernetes objects, without sending it to the API.

The following `kubectl` subcommands support this flag (non-final list):

- `apply`
- `create`
- `expose`
- `patch`
- `replace`
- `run`
- `set`

For example, we can use the `--dry-run=client` flag to create a template for our a Nginx deployment:

```bash
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml
```

The result is the following YAML output:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        resources: {}
status: {}
```

## `kubectl` API requests
If you want to see the HTTP requests kubectl sends to the Kubernetes API in detail, you can use the optional flag `--v=10`.

For example, to see the API request for creating a namespace:

```bash
kubectl create namespace test --v=10
```

The result is the following output:

```
I1119 15:47:52.822841   25474 loader.go:372] Config loaded from file:  /home/nte-jla/.kube/config
I1119 15:47:52.824692   25474 request.go:1073] Request Body: {"kind":"Namespace","apiVersion":"v1","metadata":{"name":"test2","creationTimestamp":null},"spec":{},"status":{}}
I1119 15:47:52.824891   25474 round_trippers.go:466] curl -v -XPOST  -H "User-Agent: kubectl/v1.24.3 (linux/amd64) kubernetes/aef86a9" -H "Authorization: Bearer <masked>" -H "Accept: application/json, */*" -H "Content-Type: application/json" 'https://gog-pro-lbaas-01.os.stoney-cloud.com:6443/api/v1/namespaces?fieldManager=kubectl-create&fieldValidation=Strict'
I1119 15:47:52.829841   25474 round_trippers.go:495] HTTP Trace: DNS Lookup for gog-pro-lbaas-01.os.stoney-cloud.com resolved to [{185.85.126.71 }]
I1119 15:47:52.831672   25474 round_trippers.go:510] HTTP Trace: Dial to tcp:185.85.126.71:6443 succeed
I1119 15:47:52.878279   25474 round_trippers.go:553] POST https://gog-pro-lbaas-01.os.stoney-cloud.com:6443/api/v1/namespaces?fieldManager=kubectl-create&fieldValidation=Strict 201 Created in 53 milliseconds
I1119 15:47:52.878340   25474 round_trippers.go:570] HTTP Statistics: DNSLookup 4 ms Dial 1 ms TLSHandshake 17 ms ServerProcessing 28 ms Duration 53 ms
I1119 15:47:52.878369   25474 round_trippers.go:577] Response Headers:
I1119 15:47:52.878399   25474 round_trippers.go:580]     Content-Type: application/json
I1119 15:47:52.878428   25474 round_trippers.go:580]     X-Kubernetes-Pf-Flowschema-Uid: 5bdf6f47-b545-478e-89e3-56cee0a9bfa1
I1119 15:47:52.878455   25474 round_trippers.go:580]     X-Kubernetes-Pf-Prioritylevel-Uid: f75dacbb-8f1d-4a76-8234-a205d24e39ea
I1119 15:47:52.878481   25474 round_trippers.go:580]     Content-Length: 520
I1119 15:47:52.878507   25474 round_trippers.go:580]     Date: Sat, 19 Nov 2022 14:47:52 GMT
I1119 15:47:52.878533   25474 round_trippers.go:580]     Audit-Id: 8cfa2adb-b9eb-49fc-94d2-7e1fe192c4c2
I1119 15:47:52.878559   25474 round_trippers.go:580]     Cache-Control: no-cache, private
I1119 15:47:52.878648   25474 request.go:1073] Response Body: {"kind":"Namespace","apiVersion":"v1","metadata":{"name":"test2","uid":"252882cf-a7db-4269-84c0-271381fab4d1","resourceVersion":"4295281","creationTimestamp":"2022-11-19T14:47:52Z","labels":{"kubernetes.io/metadata.name":"test2"},"managedFields":[{"manager":"kubectl-create","operation":"Update","apiVersion":"v1","time":"2022-11-19T14:47:52Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:kubernetes.io/metadata.name":{}}}}}]},"spec":{"finalizers":["kubernetes"]},"status":{"phase":"Active"}}
namespace/test2 created
```

As you can see, the output conveniently contains the corresponding `curl` commands which we could use in our own code, tools, pipelines etc.

!!! note

    If you created the deployment to see the output, you can delete it again as it’s not used anywhere else (which is also the reason why the replicas are set to `0`):

    ```bash
    kubectl delete deployment nginx
    ```
