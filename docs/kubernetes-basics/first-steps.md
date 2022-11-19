# First Steps
In this tutorial, we will interact with the Kubernetes cluster for the first time.

!!! warning

    Please make sure you completed [Setup](../setup/client-setup/index.md) before you continue with this tutorial.

## Login
To login to the stepping stone cluster head over to the [Stepping Stone Wiki](https://wiki.golog.ch/wiki/Category:Customer:_Golog_AG).
There you will find the login information for the cluster.

!!! note

    For this tutorial you can also use a local [minikube](https://minikube.sigs.k8s.io/docs/start/) cluster.

## Namespaces
As a first step on the cluster, we are going to create a new Namespace.

A Namespace is a logical design used in Kubernetes to organize and separate your applications, Deployments, Pods, Ingresses, Services, etc. on a top-level basis. 
Take a look at the [Kubernetes docs](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/). 
Authorized users inside a namespace are able to manage those resources. 
Namespace names have to be unique in your cluster.

### :octicons-tasklist-16: **Task 1**: Create a new Namespace
Create a new namespace in the tutorial environment. 
The `kubectl help` output can help you figure out the right command.

!!! note
    
    Please choose an identifying name for your Namespace, e.g. your initials or name as a prefix.

    We are going to use `<namespace>` as a placeholder for your created Namespace.

??? example "Solution"

    ```bash
    kubectl create namespace <namespace>
    ```

!!! note

    By using the following command, you can switch into another Namespace instead of specifying it for each `kubectl` command.

    **Linux/MacOS:**
    ```bash
    kubectl config set-context $(kubectl config current-context) --namespace <namespace>
    ```

    **Windows:**
    ```powershell
    kubectl config current-context
    SET KUBE_CONTEXT=[Insert output of the upper command]
    kubectl config set-context %KUBE_CONTEXT% --namespace <namespace>
    ```

    Some prefer to explicitly select the Namespace for each `kubectl` command by adding `--namespace <namespace>` or `-n <namespace>`. 
    Others prefer helper tools like `kubens`.
