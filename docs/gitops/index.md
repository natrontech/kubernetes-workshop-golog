# GitOps

## What is GitOps?

GitOps is a way to do Kubernetes application delivery. It works by using Git as a single source of truth for Kubernetes resources and everything else. With GitOps, deployments to Kubernetes are just a `git push`. This workflow allows for much easier auditing of configurations, as all changes are in source control.

## Why GitOps?

The main benefits of GitOps are:

- **Easy auditing**: All changes are in source control.
- **Easy rollbacks**: Rollbacks are just a `git revert`.
- **Easy onboarding**: Onboarding new team members is as simple as giving them access to the repository.
- **Easy CI/CD**: CI/CD pipelines are just a `git push`.

## Flux

[Flux](https://fluxcd.io/) is a tool for keeping Kubernetes clusters in sync with sources of configuration (like Git repositories). Flux works by using a Git repository as a single source of truth for Kubernetes resources. Flux monitors the repository and applies any changes to the cluster automatically.

## Argo CD

[Argo CD](https://argoproj.github.io/argo-cd/) is a declarative, GitOps continuous delivery tool for Kubernetes. Argo CD works by using a Git repository as a single source of truth for Kubernetes resources. Argo CD monitors the repository and applies any changes to the cluster automatically.

## GitLab

[GitLab](https://about.gitlab.com/) is a complete DevOps platform, delivered as a single application. It includes built-in support for GitOps workflows.
You can use Gitlab CI/CD to deploy your applications to Kubernetes or to validate your Kubernetes manifests.

### GitLab CI/CD

[GitLab CI/CD](https://docs.gitlab.com/ee/ci/) is a built-in continuous integration and continuous delivery tool that can be used to deploy your applications to Kubernetes or to validate your Kubernetes manifests.

#### GitLab Kubernetes Agent

[GitLab Kubernetes Agent](https://docs.gitlab.com/ee/user/clusters/agent/) is a tool that can help you deploy your applications to Kubernetes. It works by using a Git repository as a single source of truth for Kubernetes resources. The agent monitors the repository and applies any changes to the cluster automatically.

It is also very useful if you want to run a CI/CD pipeline against your Kubernetes cluster. The agent handels the authentication for you, so you don't have to worry about it.

#### Examples

`gitlab-ci.yml`:

```yaml
stages:
  - deploy

deploy:
    stage: deploy
    image: gitlab/gitlab-runner:latest
    script:
        - kubectl apply -f kubernetes/manifests
    only:
        - master
```

`kubernetes/manifests/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: my-app:latest
```

You can also combine the Gitlab CI with Argo CD or Flux. This way you can use the Gitlab CI to validate your Kubernetes manifests and then use Argo CD or Flux to deploy them. This is a very powerful combination.
