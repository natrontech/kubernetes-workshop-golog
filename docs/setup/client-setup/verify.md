# Verify
Verify that the `kubectl` command is available by running the following command:
```bash
kubectl version --client
```
The output should be similar to the following (1):
{ .annotate }

1.  Make sure that the version of the client is the same or higher than the version of the cluster.

```bash
Client Version: version.Info{Major:"1", Minor:"22", GitVersion:"v1.22.0", GitCommit:"cde122dc4477e5e9c5f8833d2fb01c8807a0a2b1", GitTreeState:"clean", BuildDate:"2021-06-17T20:20:38Z", GoVersion:"go1.16.5", Compiler:"gc", Platform:"linux/amd64"}
```

## Tools
You can have a look at the [Tools](tools.md) page to see what tools are available to help you manage your Kubernetes cluster.

# Next steps
The `kubectl` has many commands and options. Check them out with `kubectl --help` or `kubectl <command> --help`.
Now that you have installed `kubectl`, you can continue with the [Kubernetes Basics](../kubernetes-basics/index.md) tutorial.
