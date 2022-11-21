# StatefulSets
Stateless applications or applications with a stateful backend can be described as Deployments. However, sometimes your application has to be stateful. Examples would be an application that needs a static, non-changing hostname every time it starts or a clustered application with a strict start/stop order of its services (e.g. RabbitMQ). These features are offered by StatefulSets.

!!! note

    This tutorial does not depend on the previous ones.

## Consistent hostnames
While in normal Deployments a hash-based name of the Pods (also represented as the hostname inside the Pod) is generated, StatefulSets create Pods with preconfigured names. An example of a RabbitMQ cluster with three instances (Pods) could look like this:

```
rabbitmq-0
rabbitmq-1
rabbitmq-2
```

## Scaling
Scaling is handled differently in StatefulSets. When scaling up from 3 to 5 replicas in a Deployment, two additional Pods are started at the same time (based on the configuration). Using a StatefulSet, scaling is done serially:

Let’s use our RabbitMQ example again:

1. The StatefulSet is scaled up using: `kubectl scale deployment rabbitmq --replicas=5 --namespace <namespace>`
2. `rabbitmq-3` is started
3. As soon as Pod `rabbitmq-3` is in `Ready` state the same procedure starts for `rabbitmq-4`

When scaling down, the order is inverted. The highest-numbered Pod will be stopped first. As soon as it has finished terminating the now highest-numbered Pod is stopped. This procedure is repeated as long as the desired number of replicas has not been reached.

## Update procedure
During an update of an application with a StatefulSet the highest-numbered Pod will be the first to be updated and only after a successful start the next Pod follows.

1. Highest-numbered Pod is stopped
2. New Pod (with new image tag) is started
3. If the new Pod successfully starts, the procedure is repeated for the second highest-numbered Pod
4. And so on

If the start of a new Pod fails, the update will be interrupted so that the architecture of your application won’t break.

## Dedicated persistent volumes
A very convenient feature is that unlike a Deployment a StatefulSet makes it possible to attach a different, dedicated persistent volume to each of its Pods. This is done using a so-called VolumeClaimTemplate. This spares you from defining identical Deployments with 1 replica each but different volumes.

## Conclusion
The controllable and predictable behavior can be a perfect match for applications such as RabbitMQ or etcd, as you need unique names for such application clusters.

## :octicons-tasklist-16: **Task 1**: Create a StatefulSet
Create a file named `sts_nginx-cluster.yaml` with the following definition of a StatefulSet:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nginx-cluster
spec:
  serviceName: "nginx"
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginxinc/nginx-unprivileged:1.18-alpine
          ports:
            - containerPort: 80
              name: nginx
          resources:
            limits:
              cpu: 40m
              memory: 64Mi
            requests:
              cpu: 10m
              memory: 32Mi
```

Create the StatefulSet:

```bash
kubectl apply -f sts_nginx-cluster.yaml --namespace <namespace>
```

To watch the pods’ progress, open a second console and execute the watch command:

```bash
kubectl get pods --selector app=nginx -w --namespace <namespace>
```

!!! note

    Friendly reminder that the `kubectl get -w` command will never end unless you terminate it with `CTRL-c`.

## :octicons-tasklist-16: **Task 2**: Scale the StatefulSet
Scale the StatefulSet up:

```bash
kubectl scale statefulset nginx-cluster --replicas=3 --namespace <namespace>
```

You can again watch the pods’ progress like you did in the first task.

## :octicons-tasklist-16: **Task 3**: Update the StatefulSet
In order to update the image tag in use in a StatefulSet, you can use the `kubectl set image` command. Set the StatefulSet’s image tag to `latest`:

```bash
kubectl set image statefulset nginx-cluster nginx=docker.io/nginxinc/nginx-unprivileged:latest --namespace <namespace>
```

## :octicons-tasklist-16: **Task 4**: Rollback
Imagine you just realized that switching to the `latest` image tag was an awful idea (because it is generally not advisable). Rollback the change:

```bash
kubectl rollout undo statefulset nginx-cluster --namespace <namespace>
```

## :octicons-tasklist-16: **Task 5**: Cleanup

As with every other Kubernetes resource you can delete the StatefulSet with:

```bash
kubectl delete statefulset nginx-cluster --namespace <namespace>
```

Further information can be found in the [Kubernetes’ StatefulSet documentation](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) or this [published article](https://opensource.com/article/17/2/stateful-applications).
