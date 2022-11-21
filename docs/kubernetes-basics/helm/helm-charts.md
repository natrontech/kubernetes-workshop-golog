# Helm Charts

!!! reminder "Environment Variables"

    We are going to use some environment variables in this tutorial. Please make sure you have set them correctly.
    ```bash
    # check if the environment variables are set if not set them
    export NAMESPACE=<namespace>
    echo $NAMESPACE
    export URL=${NAMESPACE}.k8s.golog.ch
    echo $URL
    ```

In this tutorial we are going to create our very first Helm chart and deploy it.

## :octicons-tasklist-16: **Task 1**: Create a Chart
First, let’s create our chart. Open your favorite terminal and make sure you’re in the workspace for this lab, e.g. `cd ~/<workspace-kubernetes-training>`:

```bash
helm create mychart
```

You will now find a `mychart` directory with the newly created chart. It already is a valid and fully functional chart which deploys an nginx instance. Have a look at the generated files and their content. For an explanation of the files, visit the [Helm Developer Documentation](https://docs.helm.sh/developing_charts/#the-chart-file-structure). In a later section you’ll find all the information about Helm templates.

## :octicons-tasklist-16: **Task 2**: Install the Chart
Before actually deploying our generated chart, we can check the (to be) generated Kubernetes resources with the following command:

```bash
helm install --dry-run --debug --namespace $NAMESPACE myfirstrelease ./mychart
```

Finally, the following command creates a new release and deploys the application:

```bash
helm install --namespace $NAMESPACE myfirstrelease ./mychart
```

With `kubectl get pods --namespace $NAMESPACE` you should see a new Pod:

```
NAME                                     READY   STATUS    RESTARTS   AGE
myfirstrelease-mychart-4d5956b75-nd8jd   1/1     Running   0          2m21s
```

You can list the newly created Helm release with the following command:

```bash
helm list --namespace $NAMESPACE
```

## :octicons-tasklist-16: **Task 3**: Upgrade the Chart
Our freshly deployed nginx is not yet accessible from outside the Kubernetes cluster. To expose it, we have to make sure a so called ingress resource will be deployed as well.

Also make sure the application is accessible via TLS.

A look into the file `templates/ingress.yaml` reveals that the rendering of the ingress and its values is configurable through values(`values.yaml`):

```yaml
{{- if .Values.ingress.enabled -}}
{{- $fullName := include "mychart.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
{{- if and .Values.ingress.className (not (semverCompare ">=1.18-0" .Capabilities.KubeVersion.GitVersion)) }}
  {{- if not (hasKey .Values.ingress.annotations "kubernetes.io/ingress.class") }}
  {{- $_ := set .Values.ingress.annotations "kubernetes.io/ingress.class" .Values.ingress.className}}
  {{- end }}
{{- end }}
{{- if semverCompare ">=1.19-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1
{{- else if semverCompare ">=1.14-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1beta1
{{- else -}}
apiVersion: extensions/v1beta1
{{- end }}
kind: Ingress
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if and .Values.ingress.className (semverCompare ">=1.18-0" .Capabilities.KubeVersion.GitVersion) }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            {{- if and .pathType (semverCompare ">=1.18-0" $.Capabilities.KubeVersion.GitVersion) }}
            pathType: {{ .pathType }}
            {{- end }}
            backend:
              {{- if semverCompare ">=1.19-0" $.Capabilities.KubeVersion.GitVersion }}
              service:
                name: {{ $fullName }}
                port:
                  number: {{ $svcPort }}
              {{- else }}
              serviceName: {{ $fullName }}
              servicePort: {{ $svcPort }}
              {{- end }}
          {{- end }}
    {{- end }}
{{- end }}
```

Thus, we need to change this value inside our `mychart/values.yaml` file. This is also where we enable the TLS part:

!!! note

    Make sure to replace the `<url>` and `<namespace>` accordingly.

```yaml
[...]
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: nginx
    kubernetes.io/tls-acme: "true"
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/add-base-url: "true"
  hosts:
    - host: <url>
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls:
    - secretName: <namespace>-tls
      hosts:
        - <url>
[...]
```

!!! note

    Make sure to set the proper value as hostname. `<namespace>` and `<url>` will be provided by the trainer.

Before we can upgrade, make sure to delete the old ingress resource:

```bash
kubectl get ingress --namespace $NAMESPACE
kubectl delete ingress <ingress-name> --namespace $NAMESPACE
```

Apply the change by upgrading our release:

```bash
helm upgrade --namespace $NAMESPACE myfirstrelease ./mychart
```

Check whether the ingress was successfully deployed by accessing the URL.

## :octicons-tasklist-16: **Task 4**: Overwrite value using commandline param

An alternative way to set or overwrite values for charts we want to deploy is the `--set name=value` parameter. This parameter can be used when installing a chart as well as upgrading.

Update the replica count of your nginx Deployment to 2 using `--set name=value`

```bash
helm upgrade --namespace $NAMESPACE --set replicaCount=2 myfirstrelease ./mychart
```

## :octicons-tasklist-16: **Task 5**: `values.yaml`
Have a look at the `values.yaml` file in your chart and study all the possible configuration params introduced in a freshly created chart.

## :octicons-tasklist-16: **Task 6**: Uninstall the Chart
To remove an application, simply remove the Helm release with the following command:

```bash
helm uninstall --namespace $NAMESPACE myfirstrelease
```

Do this with our deployed release. With `kubectl get pods --namespace $NAMESPACE` you should no longer see your application Pod.


## Further Reading
For creating a nice README.md for your chart, have a look at [this](https://github.com/norwoodj/helm-docs) tool.
Also check out the [Artifact Hub](https://artifacthub.io/) for finding Helm charts. You can also publish your own charts there.
