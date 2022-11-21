[Helm](https://github.com/helm/helm) is a [Cloud Native Foundation](https://www.cncf.io/) project to define, install and manage applications in Kubernetes.

Helm is a Package Manager for Kubernetes

- package multiple K8s resources into a single logical deployment unit
- … but it’s not just a Package Manager

Helm is a Deployment Management for Kubernetes

- do a repeatable deployment
- manage dependencies: reuse and share
- manage multiple configurations
- update, rollback and test application deployments

## Overview
Ok, let’s start with Helm. First, you have to understand the following 3 Helm concepts: **Chart**, **Repository** and **Release**.

A Chart is a Helm package. It contains all of the resource definitions necessary to run an application, tool, or service inside of a Kubernetes cluster. Think of it like the Kubernetes equivalent of a Homebrew formula, an Apt dpkg, or a Yum RPM file.

A Repository is the place where charts can be collected and shared. It’s like Perl’s CPAN archive or the Fedora Package Database, but for Kubernetes packages.

A Release is an instance of a chart running in a Kubernetes cluster. One chart can often be installed many times in the same cluster. Each time it is installed, a new release is created. Consider a MySQL chart. If you want two databases running in your cluster, you can install that chart twice. Each one will have its own release, which will in turn have its own release name.

With these concepts in mind, we can now explain Helm like this:

!!! quote

    Helm installs charts into Kubernetes, creating a new release for each installation. To find new charts, you can search Helm chart repositories.

## Installation
This guide shows you how to install the `helm` CLI tool. `helm` can be installed either from source or from pre-built binary releases. We are going to use the pre-built releases. `helm` binaries can be found on [Helm’s release page](https://github.com/helm/helm/releases) for the usual variety of operating systems.

### :octicons-tasklist-16: **Task 1**: Install CLI

Install the CLI for your Operating System

#### Verify the installation
To verify, run the following command and check if `Version` is what you expected:

```bash
helm version
```

The output is similar to this:

```
version.BuildInfo{Version:"v3.10.0", GitCommit:"ce66412a723e4d89555dc67217607c6579ffcb21", GitTreeState:"clean", GoVersion:"go1.19.1"}
```

From here on you should be able to run the client.
