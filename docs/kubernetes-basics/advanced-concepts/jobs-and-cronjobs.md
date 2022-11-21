# Jobs and CronJobs

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    ```

Jobs are different from normal Deployments: Jobs execute a time-constrained operation and report the result as soon as they are finished; think of a batch job. To achieve this, a Job creates a Pod and runs a defined command. A Job isn’t limited to creating a single Pod, it can also create multiple Pods. When a Job is deleted, the Pods started (and stopped) by the Job are also deleted.

For example, a Job is used to ensure that a Pod is run until its completion. If a Pod fails, for example because of a Node error, the Job starts a new one. A Job can also be used to start multiple Pods in parallel.

More detailed information can be retrieved from the [Kubernetes documentation](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/).

## :octicons-tasklist-16: **Task 1**: Create a Job for a database dump
We want to create a Job that creates a postgresql database dump and stores it in a file. The Job should run once and then terminate.

Create a file called `job.yaml` with the following content:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pg-dump
spec:
  template:
    spec:
      containers:
      - name: pg-dump
        image: postgres:14.5
        command: 
        - 'bash'
        - '-eo'
        - 'pipefail'
        - '-c'
        - >
          trap "echo Backup failed; exit 0" ERR;
          FILENAME=backup-$(date +%Y-%m-%d_%H-%M-%S).sql.gz;
          echo "creating .pgpass file";
          echo "$POSTGRES_HOST:$POSTGRES_PORT:$POSTGRES_DB:$POSTGRES_USER:$POSTGRES_PASSWORD" > ~/.pgpass;
          chmod 0600 ~/.pgpass;
          echo "creating dump";
          pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER $POSTGRES_DB | gzip > $FILENAME;
          echo "";
          echo "Backup created: $FILENAME"; du -h $FILENAME;
        env:
        - name: POSTGRES_USER
          value: "<postgresql user>"
        - name: POSTGRES_HOST
          value: "<postgresql host>"
        - name: POSTGRES_PORT
          value: "5432"
        - name: POSTGRES_DB
          value: "<postgresql database>"
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: job-secret
              key: POSTGRES_PASSWORD
        resources:
          limits:
            cpu: 40m
            memory: 64Mi
          requests:
            cpu: 10m
            memory: 32Mi
      restartPolicy: Never
```

The Job uses the `postgres:14.5` image to create a database dump. The database password credentials are stored in a secret called `job-secret`. The Job is configured to run only once and then terminate.

Create the secret file `job-secret.yaml` with the following content:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: job-secret
stringData:
  POSTGRES_PASSWORD: "<postgresql password>"
```

Execute the following commands to create the Job and the secret:

```bash
kubectl apply -f job-secret.yaml --namespace $NAMESPACE
kubectl apply -f job.yaml --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 2**: Check the Job status
Check the status of the Job:

```bash
kubectl get jobs --namespace $NAMESPACE
```

The output should look like this:

```
NAME      COMPLETIONS   DURATION   AGE
pg-dump   1/1           19s        4m47
```

The Job is completed after 19 seconds. The Job created a Pod and executed the command defined in the `command` section of the Job. The output of the command is stored in the Job status.

Check the Job pod logs:

```bash
kubectl logs -f jobs/pg-dump --namespace $NAMESPACE
```

The output should look like this:

```
creating .pgpass file
creating dump

Backup created: backup-2022-11-21_09-38-50.sql.gz
4.0K    backup-2022-11-21_09-38-50.sql.gz
```

The Job created a database dump and stored it in a file called `backup-2022-11-21_09-38-50.sql.gz`.

To show all Pods belonging to a Job in a human-readable format, the following command can be used:

```bash
kubectl get pods --selector=job-name=pg-dump --output=go-template='{{range .items}}{{.metadata.name}}{{end}}' --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 3**: Create a CronJob
A CronJob is nothing else than a resource which creates a Job at a defined time, which in turn starts (as we saw in the previous section) a Pod to run a command. Typical use cases are cleanup Jobs, which tidy up old data for a running Pod, or a Job to regularly create and save a database dump as we just did during this tutorial.

The CronJob’s definition will remind you of the Deployment’s structure, or really any other control resource. There’s most importantly the `schedule` specification in [cron schedule format](https://crontab.guru/), some more things you could define and then the Job’s definition itself that is going to be created by the CronJob.

Try to create a CronJob that runs every hour at minute 0 and creates a database dump. The CronJob should be named `pg-dump-cronjob` and the Job should be named `pg-dump-job`. The Job should be created in the same namespace as the CronJob.

??? example "solution"

    ```yaml
    apiVersion: batch/v1
    kind: CronJob
    metadata:
      name: pg-dump-cronjob
    spec:
      schedule: "0 * * * *"
      jobTemplate:
        spec:
          template:
            spec:
              containers:
              - name: pg-dump
                image: postgres:14.5
                command: 
                - 'bash'
                - '-eo'
                - 'pipefail'
                - '-c'
                - >
                  trap "echo Backup failed; exit 0" ERR;
                  FILENAME=backup-$(date +%Y-%m-%d_%H-%M-%S).sql.gz;
                  echo "creating .pgpass file";
                  echo "$POSTGRES_HOST:$POSTGRES_PORT:$POSTGRES_DB:$POSTGRES_USER:$POSTGRES_PASSWORD" > ~/.pgpass;
                  chmod 0600 ~/.pgpass;
                  echo "creating dump";
                  pg_dump -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER $POSTGRES_DB | gzip > $FILENAME;
                  echo "";
                  echo "Backup created: $FILENAME"; du -h $FILENAME;
                env:
                - name: POSTGRES_USER
                  value: "<postgresql user>"
                - name: POSTGRES_HOST
                  value: "<postgresql host>"
                - name: POSTGRES_PORT
                  value: "5432"
                - name: POSTGRES_DB
                  value: "<postgresql database>"
                - name: POSTGRES_PASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: job-secret
                      key: POSTGRES_PASSWORD
                resources:
                  limits:
                    cpu: 40m
                    memory: 64Mi
                  requests:
                    cpu: 10m
                    memory: 32Mi
              restartPolicy: Never
    ```

Further information can be found in the [Kubernetes CronJob documentation](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/).
