# Database Connection
You've now created a Deployment, exposed it as a Service, and scaled it up. But how do you know if it's working? How do you know if your application is actually connecting to the database? In this module, you'll learn how to connect to your application and run some simple commands to verify that it's working.

## :octicons-tasklist-16: **Task 1**: Connect Test Webserver to Database
We've already created a Deployment. 
Now we will create another Deployment with a simple webserver that connects to the database. 
Because the application needs credentials to connect to the database, we will use a Secret to store the credentials.
Create a file called `db-secret.yaml` with the following content:

```yaml

```

!!! note

    You can generate the yaml file using the following command:

    ```bash
    kubectl create secret generic db-secret --from-literal=DB_USER=postgres --from-literal=password=mmuser_password --dry-run -o yaml > db-secret.yaml
    ```

Create a file called `db-deployment.yaml` with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: test-webserver
  name: test-webserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-webserver
  template:
    metadata:
      labels:
        app: test-webserver
    spec:
      containers:
      - image: ghcr.io/natrongmbh/kubernetes-workshop-test-webserver:latest
        name: test-webserver
        resources:
          requests:
            cpu: 10m
            memory: 16Mi
          limits:
            cpu: 20m
            memory: 32Mi
```