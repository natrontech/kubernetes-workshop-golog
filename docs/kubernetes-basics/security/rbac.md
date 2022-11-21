# RBAC

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

Role-based access control (RBAC) is a method of regulating access to computer or network resources based on the roles of individual users within your organization.
RBAC allows management of users and roles, where a role is a collection of permissions to access resources. Users can be assigned to multiple roles and permissions can be reused across multiple roles.

Read more about RBAC in the [Kubernetes documentation](https://kubernetes.io/docs/reference/access-authn-authz/rbac/).

## :octicons-tasklist-16: **Task 1**: Create a Service Account
A Service Account is an account that is used by applications to interact with the Kubernetes cluster. It is assigned a set of credentials that can be used to authenticate to the cluster. The credentials are stored as a `Secret` in the Kubernetes API.

In this task, you will create a Service Account named `app`

```bash
kubectl create serviceaccount app --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 2**: Create a Role
A Role is a collection of permissions that can be assigned to a Service Account. In this task, you will create a Role named `app` that allows the Service Account to read the `Pod` resource in the current namespace.

```bash
kubectl create role app --verb=get --verb=list --verb=watch --resource=pods --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 3**: Create a Role Binding
A Role Binding is a link between a Role and a Service Account. In this task, you will create a Role Binding that links the Role `app` to the Service Account `app`.

```bash
kubectl create rolebinding app --role=app --serviceaccount=${NAMESPACE}:app --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 4**: Verify the Role Binding
In this task, you will verify that the Role Binding is created and that the Service Account is linked to the Role.

```bash
kubectl get rolebinding app --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 5**: Create a Pod
In this task, you will create a Pod that uses the Service Account `app` to access the Kubernetes API.

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: app
  namespace: $NAMESPACE
spec:
  serviceAccountName: app
  containers:
  - name: app
    image: yauritux/busybox-curl
    command: ["sh", "-c", "sleep 3600"]
EOF
```

## :octicons-tasklist-16: **Task 6**: Verify the Pod
In this task, you will verify that the Pod is created and that it is running.

```bash
kubectl get pod app --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 7**: Access the Kubernetes API
In this task, you will access the Kubernetes API from within the Pod and verify that the Service Account has the correct permissions.

First we need the Service Account token.

```bash
# export the whole output to TOKEN
export TOKEN=$(kubectl exec app --namespace $NAMESPACE -- cat /var/run/secrets/kubernetes.io/serviceaccount/token)
```

Then we can curl the Kubernetes API over the pod's network.

```bash
kubectl exec app --namespace $NAMESPACE -- curl -ks --header "Authorization: Bearer $TOKEN" https://kubernetes.default.svc/api/v1/namespaces/$NAMESPACE/pods
```

This should return a list of Pods in the current namespace in JSON format.

## :octicons-tasklist-16: **Task 8**: Delete the Pod
In this task, you will delete the Pod.

```bash
kubectl delete pod app --namespace $NAMESPACE
```
