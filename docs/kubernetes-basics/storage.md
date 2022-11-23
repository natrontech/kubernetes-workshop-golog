# Storage

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    export URL=${NAMESPACE}.k8s.golog.ch
    echo $URL
    ```

By default, data in containers is not persistent as was the case e.g. in [Database Connection](./database-connection.md). This means that data that was written in a container is lost as soon as it does not exist anymore. We want to prevent this from happening. One possible solution to this problem is to use persistent storage.

## Request storage
Attaching persistent storage to a Pod happens in two steps. The first step includes the creation of a so-called *PersistentVolumeClaim* (PVC) in our namespace. This claim defines amongst other things what size we would like to get.

The *PersistentVolumeClaim* only represents a request but not the storage itself. It is automatically going to be bound to a PersistentVolume by Kubernetes, one that has at least the requested size. If only volumes exist that have a bigger size than was requested, one of these volumes is going to be used. The claim will automatically be updated with the new size. If there are only smaller volumes available, the claim cannot be fulfilled as long as no volume the exact same or larger size is created.

On the **stepping stone** cluster we have a `NFS` service that provides storage to the cluster. 
You can find the details of the service in the [stepping stone documentation](https://wiki.golog.ch/wiki/Category:Customer:_Golog_AG).

## NFS Storage Class
For the next steps you need to deploy a Storage Class Provider, which handles the `PVC` with the NFS Server.

A solid service is the following Project:

- [NFS Subdir External Provisioner](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner)

!!! note

    The following steps are only necessary if it was not already deployed. Also make sure you know how to use [helm](./helm/index.md)

For a simple helm deployment you can execute the following commands:

```bash
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
helm repo update
```

Create a dedicated namespace:

```bash
kubectl create namespace nfs-provisioner
```

Create the Helm deployment:

```bash
helm install nfs-subdir-external-provisioner --namespace nfs-provisioner \
  nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
  --set nfs.server=192.168.16.17 \
  --set nfs.path=/var/data/share \
  --set storageClass.name=nfs \
  --set storageClass.onDelete=true
```

## :octicons-tasklist-16: **Task 1**: Create a PersistentVolumeClaim and attach it to the Pod
For this tutorial we will create a new deployment with a simple webserver, which serves static files. We want to store the static files on a persistent volume.

!!! abstract "Details"

    For further information read the [NFS based persistent storage](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes) documentation.

Then we need to create a **PersistentVolumeClaim**.
A **persistent volume claim (PVC)** specifies the desired access mode and storage capacity. Currently, based on only these two attributes, a PVC is bound to a single PV. Once a PV is bound to a PVC, that PV is essentially tied to the PVC’s project and cannot be bound to by another PVC. There is a one-to-one mapping of PVs and PVCs. However, multiple pods in the same project can use the same PVC.

So we create a new file called `nfs-pvc.yaml` with the following content:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: webserver-pvc
spec:
  accessModes:
    - ReadWriteOnce # ReadWriteMany is also possible for NFS
  resources:
    requests:
      storage: 1Gi
  storageClassName: nfs
```

!!! note

    The `storageClassName` must match the name of the Storage Class Provider. In our case it is `nfs`. You can check the name of the Storage Class Provider with the following command:

    ```bash
    kubectl get storageclass
    ```

Afterwards, create a file called `nfs-deployment.yaml` for an classic `nginx` webserver and try to attach the persistent volume claim at the mount path `/usr/share/nginx/html` to the Pod. 

??? example "solution"

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      labels:
        app: nfs-webserver
      name: nfs-webserver
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: nfs-webserver
      template:
        metadata:
          labels:
            app: nfs-webserver
        spec:
          containers:
          - image: nginx
            name: nfs-webserver
            resources:
              requests:
                cpu: 10m
                memory: 16Mi
              limits:
                cpu: 20m
                memory: 32Mi
            volumeMounts:
            - mountPath: /usr/share/nginx/html
              name: webserver
          volumes:
          - name: webserver
            persistentVolumeClaim:
              claimName: webserver-pvc
    ```

Also create the responsible `nfs-service.yaml` and `nfs-ingress.yaml` file to expose the webserver to the outside world.

??? example "solution"

    `nfs-service.yaml`:

    ```yaml
    # nfs-service.yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: nfs-webserver
    spec:
      type: NodePort
      ports:
      - port: 80
        targetPort: 80
      selector:
        app: nfs-webserver
    ```

    `nfs-ingress.yaml`:

    ```bash
    kubectl create --dry-run=client --namespace $NAMESPACE -o yaml -f - <<EOF >> nfs-ingress.yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: nfs-webserver
      annotations:
        kubernetes.io/ingress.class: nginx
        kubernetes.io/tls-acme: "true"
        cert-manager.io/cluster-issuer: letsencrypt-prod
        nginx.ingress.kubernetes.io/add-base-url: "true"
    spec:
      tls:
      - hosts:
        - $URL
        secretName: ${NAMESPACE}-tls
      rules:
      - host: $URL
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nfs-webserver
                port:
                  number: 80
    EOF
    ```

!!! note "Ingress"
    
    Before you can deploy the Ingress, make sure you delete the old one. Otherwise you will get an error message.

    ```bash
    kubectl get ingress --namespace $NAMESPACE
    kubectl delete ingress <ingress> --namespace $NAMESPACE
    ```

Now you can deploy everything in order to test it.


```bash
kubectl apply -f nfs-pv.yaml --namespace $NAMESPACE
kubectl apply -f nfs-pvc.yaml --namespace $NAMESPACE
kubectl apply -f nfs-deployment.yaml --namespace $NAMESPACE
kubectl apply -f nfs-service.yaml --namespace $NAMESPACE
kubectl apply -f nfs-ingress.yaml --namespace $NAMESPACE
```

Describe the Persistent Volume Claim to see if it was created successfully.

## :octicons-tasklist-16: **Task 2**: Adding a file to the persistent volume
When we now visit the webserver, we will see that we get a `403 Forbidden` error. This is because the directory is empty and we need to add a file to it. We can do this by opening a shell in the Pod and adding a file to the directory.

There are multiple ways to do this. One way is to use the `kubectl exec` command. This command allows us to execute a command in a running container. We can use this to open a shell in the container.

```bash
kubectl exec -it nfs-webserver-<pod-id> --namespace $NAMESPACE -- /bin/bash
```

Because there is no editor installed in the container, we simply use `echo` to create a file.

```bash
echo "<h1>Hello World</h1>" > /usr/share/nginx/html/index.html
```

!!! note

    You can also execute the `echo` command directly in the `kubectl exec` command.

    ```bash
    kubectl exec -it nfs-webserver-<pod-id> --namespace $NAMESPACE -- /bin/bash -c "echo '<h1>Hello World</h1>' > /usr/share/nginx/html/index.html"
    ```

Now we can visit the webserver and see that the file was created successfully.

**Another way** to do this is to use the `kubectl cp` command. This command allows us to copy files from and to a container. We can use this to copy a file from our local machine to the container.

First, we need to create a `index.html` file on our local machine.

```bash
echo "<h1>Hello World</h1>" > index.html
```

Then we can copy the file to the container.

```bash
kubectl cp index.html nfs-webserver-<pod-id>:/usr/share/nginx/html/index.html --namespace $NAMESPACE
```

Now we can visit the webserver and see that the file was created successfully.


## :octicons-tasklist-16: **Task 3**: Deleting the Pod
Now we can delete the Pod and see that the file is still there. This is because the file is stored on the NFS server and not in the Pod.

```bash
export POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l app=nfs-webserver -o jsonpath="{.items[0].metadata.name}")
kubectl delete pod $POD_NAME --namespace $NAMESPACE
```

Now we can visit the webserver and see that the file is still there.
