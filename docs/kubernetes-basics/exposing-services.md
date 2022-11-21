# Exposing Services

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    export URL=${NAMESPACE}.k8s.golog.ch
    echo $URL
    ```

In this module, you'll learn how to expose an application to the outside world.

## :octicons-tasklist-16: **Task 1**: Create a NodePort Service with an Ingress
The command `kubectl apply -f deployment.yaml `from the last tutorial creates a Deployment but no Service. A Kubernetes Service is an abstract way to expose an application running on a set of Pods as a network service. For some parts of your application (for example, frontends) you may want to expose a Service to an external IP address which is outside your cluster.

Kubernetes `ServiceTypes` allow you to specify what kind of Service you want. The default is `ClusterIP`.

`Type` values and their behaviors are:

- `ClusterIP`: Exposes the Service on a cluster-internal IP. Choosing this value only makes the Service reachable from within the cluster. This is the default ServiceType.
- `NodePort`: Exposes the Service on each Node’s IP at a static port (the NodePort). A ClusterIP Service, to which the NodePort Service routes, is automatically created. You’ll be able to contact the NodePort Service from outside the cluster, by requesting <NodeIP>:<NodePort>.
- `LoadBalancer`: Exposes the Service externally using a cloud provider’s load balancer. NodePort and ClusterIP Services, to which the external load balancer routes, are automatically created.
- `ExternalName`: Maps the Service to the contents of the externalName field (e.g. foo.bar.example.com), by returning a CNAME record with its value. No proxying of any kind is set up.

You can also use Ingress to expose your Service. Ingress is not a Service type, but it acts as the entry point for your cluster. [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource. An Ingress may be configured to give Services externally reachable URLs, load balance traffic, terminate SSL / TLS, and offer name-based virtual hosting. An Ingress controller is responsible for fulfilling the route, usually with a load balancer, though it may also configure your edge router or additional frontends to help handle the traffic.

In order to create an Ingress, we first need to create a Service of type [ClusterIP](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types) . We’re going to do this with the command `kubectl expose`:

```bash
kubectl expose deployment/test-webserver --name=test-webserver --port=8080 --target-port=8080 --type=NodePort --namespace $NAMESPACE
```

Let’s have a more detailed look at our Service:

```bash
kubectl get service test-webserver --namespace $NAMESPACE
```

The output should look like this:

```bash
NAME                TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
test-webserver      NodePort   10.97.53.32   <none>        8080:32329/TCP   4s
```

!!! note

    Service IP (CLUSTER-IP) addresses stay the same for the duration of the Service’s lifespan.

By executing the following command:

```bash
kubectl get service test-webserver -o yaml --namespace $NAMESPACE
```

You get additional information:

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: test-webserver
  name: test-webserver
  namespace: test-ns
  resourceVersion: "4270474"
spec:
  clusterIP: 10.97.53.32
  clusterIPs:
  - 10.97.53.32
  externalTrafficPolicy: Cluster
  internalTrafficPolicy: Cluster
  ipFamilies:
  - IPv4
  ipFamilyPolicy: SingleStack
  ports:
  - nodePort: 32329
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: test-webserver
  sessionAffinity: None
  type: NodePort
status:
  loadBalancer: {}
```

The Service’s `selector` defines which Pods are being used as Endpoints. This happens based on labels. Look at the configuration of Service and Pod in order to find out what maps to what:

```bash
kubectl get service test-webserver -o yaml --namespace $NAMESPACE
```

```yaml
...
  selector:
    app: test-webserver
...
```

With the following command you get details from the Pod:

!!! note

    First, get all Pod names from your namespace with (`kubectl get pods --namespace $NAMESPACE`) and then replace <pod> in the following command. If you have installed and configured the bash completion, you can also press the TAB key for autocompletion of the Pod’s name.

```bash
export POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l "app=test-webserver" -o jsonpath="{.items[0].metadata.name}")
kubectl get pod $POD_NAME -o yaml --namespace $NAMESPACE
```

Let’s have a look at the label section of the Pod and verify that the Service selector matches the Pod’s labels:

```yaml
...
  labels:
    app: test-webserver
...
```

This link between Service and Pod can also be displayed in an easier fashion with the kubectl describe command:

```bash
kubectl describe service test-webserver --namespace $NAMESPACE
```

```
Name:                     test-webserver
Namespace:                test-ns
Labels:                   app=test-webserver
Annotations:              <none>
Selector:                 app=test-webserver
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.97.53.32
IPs:                      10.97.53.32
Port:                     <unset>  8080/TCP
TargetPort:               8080/TCP
NodePort:                 <unset>  32329/TCP
Endpoints:                <none>
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

The `Endpoints` show the IP addresses of all currently matched Pods.

With the NodePort Service ready, we can now create the Ingress resource.

In order to create the Ingress resource, we first need to create the file `ingress.yaml` and change the host entry to match your environment:

```bash
kubectl create --dry-run=client --namespace $NAMESPACE -o yaml -f - <<EOF >> ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: test-webserver
  annotations:
    kubernetes.io/ingress.class: nginx
    kubernetes.io/tls-acme: "true"
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/add-base-url: "true"
spec:
  tls:
  - hosts:
    - $URL
    secretName: ${URL}-test-webserver-tls
  rules:
  - host: $URL
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: test-webserver
            port:
              number: 8080
EOF
```

As you see in the resource definition at `spec.rules[0].http.paths[0].backend.service.name` we use the previously created `test-webserver` NodePort Service.

Let’s create the Ingress resource with:

```bash
kubectl apply -f ingress.yaml --namespace $NAMESPACE
```

Get the hostname of the Ingress resource:

```bash
kubectl get ingress test-webserver --namespace $NAMESPACE
```

Afterwards, we are able to access our freshly created Ingress at `https://<namespace>.k8s.golog.ch`

## :octicons-tasklist-16: **Task 2**: For fast learners
Have a closer look at the resources created in your namespace $NAMESPACE with the following commands and try to understand them:

```bash
kubectl describe namespace $NAMESPACE
```

```bash
kubectl get all --namespace $NAMESPACE
```

```bash
kubectl describe <resource> <name> --namespace $NAMESPACE
```

```bash
kubectl get <resource> <name> -o yaml --namespace $NAMESPACE
```
