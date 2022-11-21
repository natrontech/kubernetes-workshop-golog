# Configmaps

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

Similar to environment variables, ConfigsMaps allow you to separate the configuration for an application from the image. Pods can access those variables at runtime which allows maximum portability for applications running in containers. In this lab, you will learn how to create and use ConfigMaps.

A ConfigMap can be created using the `kubectl create configmap` command as follows:

```bash
kubectl create configmap <name> <data-source> --namespace $NAMESPACE
```

Where the `<data-source>` can be a file, directory, or command line input.

## :octicons-tasklist-16: **Task 1**: Create a ConfigMap for Java properties
A classic example for ConfigMaps are properties files of Java applications which can’t be configured with environment variables.

First, create a file called `java.properties` with the following content:

```
JAVA_OPTS=-Xmx512m
key=value
key2=value2
```

Now you can create a ConfigMap based on that file:

```bash
kubectl create configmap javaconfiguration --from-file=./java.properties --namespace $NAMESPACE
```

Verify that the ConfigMap was created successfully:

```bash
kubectl get configmaps --namespace $NAMESPACE
```

The output should look like this:

```
NAME               DATA   AGE
javaconfiguration   1      2m
```

Have a look at its content:

```bash
kubectl get configmap javaconfiguration -o yaml --namespace $NAMESPACE
```

The output should look like this:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: javaconfiguration
data:
  java.properties: |
    JAVA_OPTS=-Xmx512m
    key=value
    key2=value2
```

## :octicons-tasklist-16: **Task 2**: Create a Pod that uses the ConfigMap
Next, we want to make a ConfigMap accessible for a container. There are basically the following possibilities to achieve [this](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/):

- ConfigMap properties as environment variables in a Deployment
- Command line arguments via environment variables
- Mounted as volumes in the container

In this example, we want the file to be mounted as a volume inside the container.

Basically, a Deployment has to be extended with the following config:

```yaml
      ...
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      ...
      volumes:
      - configMap:
          defaultMode: 420
          name: javaconfiguration
        name: config-volume
      ...
```

The `volumeMounts` section defines the mount point inside the container. The `volumes` section defines the volume that should be mounted. The `name` property of the volume has to match the `name` property of the `volumeMounts` section.

Create a file called `java-deployment.yaml` with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: spring-boot-example
  name: spring-boot-example
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: spring-boot-example
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: spring-boot-example
    spec:
      containers:
        - image: appuio/example-spring-boot
          imagePullPolicy: Always
          name: example-spring-boot
          resources: 
            limits:
              cpu: 1
              memory: 768Mi
            requests:
              cpu: 20m
              memory: 32Mi
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - mountPath: /etc/config
              name: config-volume
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
        - configMap:
            defaultMode: 420
            name: javaconfiguration
          name: config-volume
```

Deploy it:

```bash
kubectl apply -f java-deployment.yaml --namespace $NAMESPACE
```

This means that the container should now be able to access the ConfigMap’s content in `/etc/config/java.properties`. Let’s check:

```bash
export POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l "app=spring-boot-example" -o jsonpath="{.items[0].metadata.name}")
kubectl exec -it $POD_NAME --namespace $NAMESPACE -- cat /etc/config/java.properties
```

!!! note

    On Windows, you can use Git Bash with `winpty kubectl exec -it <pod> --namespace $NAMESPACE -- cat //etc/config/java.properties`.

The output should look like this:

```
JAVA_OPTS=-Xmx512m
key=value
key2=value2
```

Like this, the property file can be read and used by the application inside the container. The image stays portable to other environments.

## :octicons-tasklist-16: **Task 3**: Create a ConfigMap for environment variables

In the previous task, we created a ConfigMap for a Java properties file. Now we want to create a ConfigMap for environment variables.

Use a ConfigMap to define the environment variables `JAVA_OPTS` and `JAVA_TOOL_OPTIONS`. 
You can refer to the [official documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#define-container-environment-variables-using-configmap-data) for more information.

??? example "solutions"

    ```bash
    kubectl create configmap javaenv --from-literal=JAVA_OPTS=-Xmx512m --from-literal=JAVA_TOOL_OPTIONS=-Dfile.encoding=UTF-8 --namespace $NAMESPACE
    ```

    Update the `java-deployment.yaml` file with the following content:

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      labels:
        app: spring-boot-example
      name: spring-boot-example
    spec:
      progressDeadlineSeconds: 600
      replicas: 1
      revisionHistoryLimit: 10
      selector:
        matchLabels:
          app: spring-boot-example
      strategy:
        rollingUpdate:
          maxSurge: 25%
          maxUnavailable: 25%
        type: RollingUpdate
      template:
        metadata:
          labels:
            app: spring-boot-example
        spec:
          containers:
            - image: appuio/example-spring-boot
              imagePullPolicy: Always
              name: example-spring-boot
              resources: 
                limits:
                  cpu: 1
                  memory: 768Mi
                requests:
                  cpu: 20m
                  memory: 32Mi
              terminationMessagePath: /dev/termination-log
              terminationMessagePolicy: File
              volumeMounts:
                - mountPath: /etc/config
                  name: config-volume
              envFrom:
                - configMapRef:
                    name: javaenv
          dnsPolicy: ClusterFirst
          restartPolicy: Always
          schedulerName: default-scheduler
          securityContext: {}
          terminationGracePeriodSeconds: 30
          volumes:
            - configMap:
                defaultMode: 420
                name: javaconfiguration
              name: config-volume
    ```

    Make sure you delete the old deployment first:

    ```bash
    kubectl delete deployment spring-boot-example --namespace $NAMESPACE
    ```

    Then, create the new deployment:

    ```bash
    kubectl create -f java-deployment.yaml --namespace $NAMESPACE
    ```

    Check the environment variables of the container:

    ```bash
    export POD_NAME=$(kubectl get pods --namespace $NAMESPACE -l "app=spring-boot-example" -o jsonpath="{.items[0].metadata.name}")
    kubectl exec -it $POD_NAME --namespace $NAMESPACE -- env
    ```

    The output should look like this:

    ```
    ...
    JAVA_OPTS=-Xmx512m
    JAVA_TOOL_OPTIONS=-Dfile.encoding=UTF-8
    ...
    ```