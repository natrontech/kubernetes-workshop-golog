# Security Contexts

In the concept of security context for a pod or container, there are severals thing to consider:

- Access control
- SElinux
- Running privileged or unprivileged workload
- Linux capabilities
- AppArmor
- Seccomp

In this tutorial you will learn where to configure and how to use some of these types.

## :octicons-tasklist-16: **Task 1**: Access Control
Create a new pod by using this example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
  volumes:
  - name: sec-ctx-vol
    emptyDir: {}
  containers:
  - name: sec-ctx-demo
    image: busybox:1.28
    command: [ "sh", "-c", "sleep 1h" ]
    volumeMounts:
    - name: sec-ctx-vol
      mountPath: /data/demo
    securityContext:
      allowPrivilegeEscalation: false
```

You can see the different value entries in the ‘securityContext’ section, let’s figure how what do they do. 
So create the pod and connect into the shell:

```bash
kubectl exec -it security-context-demo --namespace <namespace> -- sh
```

In the container run `ps` to get a list of all running processes. The output shows, that the processes are running with the user `1000`, which is the value from `runAsUser`:

```
/ $ ps
PID   USER     TIME  COMMAND
    1 1000      0:00 sleep 1h
    6 1000      0:00 sh
   12 1000      0:00 ps
```

Now navigate to the directory `/data` and list the content. As you can see the `emptyDir` has been mounted with the group ID of `2000`, which is the value of the `fsGroup` field.

```
/data $ ls -lah
total 0
drwxr-xr-x    3 root     root          18 Nov 21 13:12 .
drwxr-xr-x    1 root     root          63 Nov 21 13:12 ..
drwxrwsrwx    2 root     2000           6 Nov 21 13:12 demo
```

Go into the dir `demo` and create a file:

```
/data $ cd demo/
/data/demo $ echo hello > demofile
/data/demo $ ls -lah
total 4
drwxrwsrwx    2 root     2000          22 Nov 21 13:15 .
drwxr-xr-x    3 root     root          18 Nov 21 13:12 ..
-rw-r--r--    1 1000     2000           6 Nov 21 13:15 demofile
```

List the content with `ls -lah` again and see, that `demofile` has the group ID `2000`, which is the value `fsGroup` as well.

Run the last command `id` here and check the output:

```
/data/demo $ id
uid=1000 gid=3000 groups=2000
```

The shown group ID of the user is `3000`, from the field `runAsGroup`. If the field would be empty the user would have 0 (root) and every process would be able to go with files which are owned by the root (0) group.

```bash
/data/demo $ exit
```

Check out the documentation at kubernetes.io for more information about [Security Context](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).
