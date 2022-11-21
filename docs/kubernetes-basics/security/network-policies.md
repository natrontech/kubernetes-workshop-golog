# Network Policies

One CNI function is the ability to enforce network policies and implement an in-cluster zero-trust container strategy. Network policies are a default Kubernetes object for controlling network traffic, but a CNI such as [Cilium](https://cilium.io/) or [Calico](https://www.tigera.io/project-calico/) is required to enforce them. We will demonstrate traffic blocking with our simple app.

!!! warning

    This section requires a CNI that supports network policies. The Flannel CNI does not support network policies which is currently the default CNI for a stepping stone cluster.

!!! note

    If you are not yet familiar with Kubernetes Network Policies we suggest going to the [Kubernetes Documentation](https://kubernetes.io/docs/concepts/services-networking/network-policies/).

## :octicons-tasklist-16: **Task 1**: Deploy a simple frontend/backend application
First we need a simple application to show the effects on Kubernetes network policies. Let’s have a look at the following resource definitions:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  labels:
    app: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend-container
        image: docker.io/byrnedo/alpine-curl:0.1.8
        imagePullPolicy: IfNotPresent
        command: [ "/bin/ash", "-c", "sleep 1000000000" ]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: not-frontend
  labels:
    app: not-frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: not-frontend
  template:
    metadata:
      labels:
        app: not-frontend
    spec:
      containers:
      - name: not-frontend-container
        image: docker.io/byrnedo/alpine-curl:0.1.8
        imagePullPolicy: IfNotPresent
        command: [ "/bin/ash", "-c", "sleep 1000000000" ]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend-container
        env:
        - name: PORT
          value: "8080"
        ports:
        - containerPort: 8080
        image: docker.io/cilium/json-mock:1.2
        imagePullPolicy: IfNotPresent
---
apiVersion: v1
kind: Service
metadata:
  name: backend
  labels:
    app: backend
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - name: http
    port: 8080
```

The application consists of two client deployments (`frontend` and `not-frontend`) and one backend deployment (`backend`). We are going to send requests from the frontend and not-frontend pods to the backend pod.

Create a file `policies-app.yaml` with the above content.

Deploy the app in a new namespace:

```bash
export NAMESPACE=<username>-policies
kubectl create namespace $NAMESPACE
kubectl apply -f policies-app.yaml --namespace $NAMESPACE
```

this gives you the following output:

```
deployment.apps/frontend created
deployment.apps/not-frontend created
deployment.apps/backend created
service/backend created
```

Verify with the following command that everything is up and running:

```bash
kubectl get pods --namespace $NAMESPACE
```

Let us make life a bit easier by storing the pods name into an environment variable so we can reuse it later again:

```bash
export FRONTEND=$(kubectl get pods -l app=frontend -o jsonpath='{.items[0].metadata.name}' --namespace $NAMESPACE)
echo ${FRONTEND}
export NOT_FRONTEND=$(kubectl get pods -l app=not-frontend -o jsonpath='{.items[0].metadata.name}' --namespace $NAMESPACE)
echo ${NOT_FRONTEND}
```

## :octicons-tasklist-16: **Task 2**: Verify that the frontend can access the backend
Now we generate some traffic as a baseline test.

```bash
kubectl exec -ti ${FRONTEND} --namespace $NAMESPACE -- curl -I --connect-timeout 5 backend:8080
```

and

```bash
kubectl exec -ti ${NOT_FRONTEND} --namespace $NAMESPACE -- curl -I --connect-timeout 5 backend:8080
```

This will execute a simple `curl` call from the `frontend` and `not-frondend` application to the backend application:

```
HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Mon, 21 Nov 2022 13:00:59 GMT
Connection: keep-alive

HTTP/1.1 200 OK
X-Powered-By: Express
Vary: Origin, Accept-Encoding
Access-Control-Allow-Credentials: true
Accept-Ranges: bytes
Cache-Control: public, max-age=0
Last-Modified: Sat, 26 Oct 1985 08:15:00 GMT
ETag: W/"83d-7438674ba0"
Content-Type: text/html; charset=UTF-8
Content-Length: 2109
Date: Mon, 21 Nov 2022 13:01:18 GMT
Connection: keep-alive
```

and we see, both applications can connect to the `backend` application.

Until now ingress and egress policy enforcement are still disabled on all of our pods because no network policy has been imported yet selecting any of the pods. Let us change this.

## :octicons-tasklist-16: **Task 3**: Deny traffic with a network policy

We block traffic by applying a network policy. Create a file `policies-backend-ingress-deny.yaml` with the following content:

```yaml
---
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: backend-ingress-deny
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
```

The policy will deny all ingress traffic as it is of type Ingress but specifies no allow rule, and will be applied to all pods with the `app=backend` label thanks to the podSelector.

Ok, then let’s create the policy with:

```bash
kubectl apply -f policies-backend-ingress-deny.yaml --namespace $NAMESPACE
```

and verify that the policy has been created:

```bash
kubectl get networkpolicies --namespace $NAMESPACE
```

which gives you the following output:
```
NAME                   POD-SELECTOR   AGE
backend-ingress-deny   app=backend    4s
```

## :octicons-tasklist-16: **Task 4**: Verify that the frontend can no longer access the backend
We can now execute the connectivity check again:

```bash
kubectl exec -ti ${FRONTEND} --namespace $NAMESPACE -- curl -I --connect-timeout 5 backend:8080
```

and

```bash
kubectl exec -ti ${NOT_FRONTEND} --namespace $NAMESPACE -- curl -I --connect-timeout 5 backend:8080
```

but this time you see that the `frontend` and `not-frontend` application cannot connect anymore to the backend:

```
# Frontend
curl: (28) Connection timed out after 5001 milliseconds
command terminated with exit code 28
# Not Frontend
curl: (28) Connection timed out after 5001 milliseconds
command terminated with exit code 28
```
