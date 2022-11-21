# ResourceQuotas and LimitRanges

In this tutorial, we are going to look at ResourceQuotas and LimitRanges. As Kubernetes users, we are most certainly going to encounter the limiting effects that ResourceQuotas and LimitRanges impose.

!!! warning

    For this lab to work it is vital that you create and use the namespace `<username>-quota`!

## ResourceQuotas
ResourceQuotas among other things limit the amount of resources Pods can use in a Namespace. They can also be used to limit the total number of a certain resource type in a Namespace. In more detail, there are these kinds of quotas:

- *Compute ResourceQuotas* can be used to limit the amount of memory and CPU
- *Storage ResourceQuotas* can be used to limit the total amount of storage and the number of PersistentVolumeClaims, generally or specific to a StorageClass
- *Object count quotas* can be used to limit the number of a certain resource type such as Services, Pods or Secrets

Defining ResourceQuotas makes sense when the cluster administrators want to have better control over consumed resources. A typical use case are public offerings where users pay for a certain guaranteed amount of resources which must not be exceeded.

In order to check for defined quotas in your Namespace, simply see if there are any of type ResourceQuota:

```bash
kubectl get resourcequota --namespace <namespace>
```

To show in detail what kinds of limits the quota imposes:

```bash
kubectl describe resourcequota <quota-name> --namespace <namespace>
```

For more details, have look at [Kubernetes’ documentation about resource quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/#requests-vs-limits).

## Requests and limits
As we’ve already seen, compute ResourceQuotas limit the amount of memory and CPU we can use in a Namespace. Only defining a ResourceQuota, however is not going to have an effect on Pods that don’t define the amount of resources they want to use. This is where the concept of limits and requests comes into play.

Limits and requests on a Pod, or rather on a container in a Pod, define how much memory and CPU this container wants to consume at least (request) and at most (limit). Requests mean that the container will be guaranteed to get at least this amount of resources, limits represent the upper boundary which cannot be crossed. Defining these values helps Kubernetes in determining on which Node to schedule the Pod because it knows how many resources should be available for it.

!!! note

    Containers using more CPU time than what their limit allows will be throttled. Containers using more memory than what they are allowed to use will be killed.

Defining limits and requests on a Pod that has one container looks like this:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: lr-demo
  namespace: lr-example
spec:
  containers:
  - name: lr-demo-ctr
    image: docker.io/nginxinc/nginx-unprivileged:1.18-alpine
    resources:
      limits:
        memory: "200Mi"
        cpu: "700m"
      requests:
        memory: "200Mi"
        cpu: "700m"
```

You can see the familiar binary unit “Mi” is used for the memory value. Other binary (“Gi”, “Ki”, …) or decimal units (“M”, “G”, “K”, …) can be used as well.

The CPU value is denoted as “m”. “m” stands for millicpu or sometimes also referred to as *millicores* where "`1000m`" is equal to one core/vCPU/hyperthread.

### Quality of service
Setting limits and requests on containers has yet another effect: It might change the Pod’s Quality of Service class. There are three such QoS classes:

- *Guaranteed*
- *Burstable*
- *BestEffort*

The Guaranteed QoS class is applied to Pods that define both limits and requests for both memory and CPU resources on all their containers. The most important part is that each request has the same value as the limit. Pods that belong to this QoS class will never be killed by the scheduler because of resources running out on a Node.

!!! note

    If a container only defines its limits, Kubernetes automatically assigns a request that matches the limit.

The Burstable QoS class means that limits and requests on a container are set, but they are different. It is enough to define limits and requests on one container of a Pod even though there might be more, and it also only has to define limits and requests on memory or CPU, not necessarily both.

The BestEffort QoS class applies to Pods that do not define any limits and requests at all on any containers. As its class name suggests, these are the kinds of Pods that will be killed by the scheduler first if a Node runs out of memory or CPU. As you might have already guessed by now, if there are no BestEffort QoS Pods, the scheduler will begin to kill Pods belonging to the class of Burstable. A Node hosting only Pods of class Guaranteed will (theoretically) never run out of resources.

For more examples have a look at the [Kubernetes documentation about Quality of Service](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/).

## LimitRanges
As you now know what limits and requests are, we can come back to the statement made above:

!!! quote

    As we’ve already seen, compute ResourceQuotas limit the amount of memory and CPU we can use in a Namespace. Only defining a ResourceQuota, however is not going to have an effect on Pods that don’t define the amount of resources they want to use. This is where the concept of limits and requests comes into play.

So, if a cluster administrator wanted to make sure that every Pod in the cluster counted against the compute ResourceQuota, the administrator would have to have a way of defining some kind of default limits and requests that were applied if none were defined in the containers. This is exactly what LimitRanges are for.

Quoting the [Kubernetes documentation](https://kubernetes.io/docs/concepts/policy/limit-range/), LimitRanges can be used to:

- Enforce minimum and maximum compute resource usage per Pod or container in a Namespace
- Enforce minimum and maximum storage request per PersistentVolumeClaim in a Namespace
- Enforce a ratio between request and limit for a resource in a Namespace
- Set default request/limit for compute resources in a Namespace and automatically inject them to containers at runtime

If for example a container did not define any requests or limits and there was a LimitRange defining the default values, these default values would be used when deploying said container. However, as soon as limits or requests were defined, the default values would no longer be applied.

The possibility of enforcing minimum and maximum resources and defining ResourceQuotas per Namespace allows for many combinations of resource control.

## :octicons-tasklist-16: **Task 1**: Play with ResourceQuotas and LimitRanges
In this task, we will play around with ResourceQuotas and LimitRanges to get a better understanding of how they work.

### Create a Namespace
First, we will create a new Namespace for our experiments:

```bash
export NAMESPACE=<username>-quota
kubectl create namespace $NAMESPACE
```

### Create a ResourceQuota
Next, we will create a ResourceQuota for our Namespace:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ResourceQuota
metadata:
  name: lr-quota
  namespace: $NAMESPACE
spec:
  hard:
    limits.cpu: "1"
    limits.memory: 1Gi
    requests.cpu: "1"
    requests.memory: 1Gi
EOF
```

### Create a LimitRange
Now, we will create a LimitRange for our Namespace:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: LimitRange
metadata:
  name: lr-range
  namespace: $NAMESPACE
spec:
  limits:
  - default:
      cpu: 500m
      memory: 200Mi
    defaultRequest:
      cpu: 500m
      memory: 200Mi
    type: Container
EOF
```

### Create a Pod
We create a Pod that does uses a lot of resources:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: lr-demo
  namespace: $NAMESPACE
spec:
  containers:
  - name: lr-demo-ctr
    image: bretfisher/stress:2cpu1024m
EOF
```

Do you see any problems with the Pod? If not, let’s have a look at the events:

```bash
kubectl get pods -n $NAMESPACE
kubectl get events -n $NAMESPACE
```

See the error message? It says that the Pod is in a state of `OOMKilled` which stands for “Out of Memory Killed”. This is because the Pod requested more memory than the ResourceQuota allows. The Pod was killed by the scheduler.
