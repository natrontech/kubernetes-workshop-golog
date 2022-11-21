# Init Containers

!!! warning "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

A Pod can have multiple containers running apps within it, but it can also have one or more init containers, which are run before the app container is started.

Init containers are exactly like regular containers, except:

- Init containers always run to completion.
- Each init container must complete successfully before the next one starts.

Check [Init Containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) from the Kubernetes documentation for more details.

## :octicons-tasklist-16: **Task 1**: Create an init container
We want to create a Pod that runs an init container that creates a file and then runs a container that reads the file.

Create a file called `init-container.yaml` with the following content:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-demo
spec:
  containers:
  - name: nginx
    image: nginx
    resources:
      requests:
        cpu: 10m
        memory: 16Mi
      limits:
        cpu: 20m
        memory: 32Mi
    volumeMounts:
    - name: workdir
      mountPath: /usr/share/nginx/html
  initContainers:
  - name: init
    image: busybox
    command:
    - 'sh'
    - '-c'
    - 'echo "<h1>Hello World</h1>" > /work-dir/index.html'
    volumeMounts:
    - name: workdir
      mountPath: /work-dir
  volumes:
  - name: workdir
    emptyDir: {}
```

Apply the file:

```bash
kubectl apply -f init-container.yaml --namespace $NAMESPACE
```

Check the status of the Pod:

```bash
kubectl get pod init-demo --namespace $NAMESPACE
```

The output should look like this:

```bash
NAME        READY   STATUS     RESTARTS   AGE
init-demo   0/1     Init:0/1   0          1m
```

The Pod is in the `Init:0/1` state, which means that the init container is running.

Check the `/usr/share/nginx/html/index.html` file in the `nginx` container:

```bash
kubectl exec -it init-demo --namespace $NAMESPACE -- cat /usr/share/nginx/html/index.html
```

The output should look like this:

```
Defaulted container "nginx" out of: nginx, init (init)
<h1>Hello World</h1>
```

!!! question "for fast learners"
    
    Try to expose the Pod with a Service and an Ingress.

## :octicons-tasklist-16: **Task 2**: Create a Init Container that checks database connectivity
We want to create a Deployment that runs an init container that checks database connectivity and then starts the main container.

Create a file called `init-deployment.yaml` with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: init-deployment
  name: init-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: init-deployment
  template:
    metadata:
      labels:
        app: init-deployment
    spec:
      initContainers:
      - name: pg-isready
        image: postgres:14.5
        command: 
          - 'sh'
          - '-c'
          - |
            until pg_isready -h <postgresql host> -p 5432; do
              echo "waiting for database to start"
              sleep 2
            done
          - 'echo'
          - 'Database is ready'
      containers:
      - image: ghcr.io/natrongmbh/kubernetes-workshop-golog-postgresql-webserver:latest
        name: init-deployment
        resources:
          requests:
            cpu: 10m
            memory: 16Mi
          limits:
            cpu: 20m
            memory: 32Mi
        envFrom:
        - secretRef:
            name: db-secret
```

!!! note

    Make sure to replace `<postgresql host>` with the host of your PostgreSQL database.

Apply the file:

```bash
kubectl apply -f init-deployment.yaml --namespace $NAMESPACE
```

Check the status of the Deployment:

```bash
kubectl get deployment init-deployment --namespace $NAMESPACE
```

The output should look like this:

```
NAME                                        READY   STATUS      RESTARTS   AGE
init-deployment-5567dc778c-ns27h            0/1     Init:0/1    0          3s
```

The Deployment is in the `Init:0/1` state, which means that the init container is running.

Check the logs of the init container:

```bash
export POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l "app=init-deployment" -o jsonpath="{.items[0].metadata.name}")
kubectl logs $POD_NAME -c pg-isready --namespace $NAMESPACE
```
