# Introduction
In this lab, we introduced the core concepts of Kubernetes.  
All instructions and resources used in this lab are for quick overview only, not detailed instructions.   
Please check the official documentation for more details. 

- [Kubernetes Docs](https://kubernetes.io/docs/home/)

## Core idea 
With open-source software Kubernetes, you get a platform that can deploy your application in containers and run it simultaneously.  
For this reason, Kubernetes is also known as a container platform, or to use the term container as a service (CaaS).  
Depending on the configuration, the term Platform as a Service (PaaS) also applies. 

### Container engine 
The underlying container engine of Kubernetes is mostly Docker.  
There are other container engines that can be used with Kubernetes, such as CRI-O.

Docker was originally created to help developers test their applications in their continuous integration environments. 
Nowadays, system admins also use it. CRI-O doesn’t exist as long as Docker does. 
It is a “lightweight container runtime for Kubernetes” and is fully [OCI-compliant](https://github.com/opencontainers/runtime-spec).

## Overview
Kubernetes consists of control plane and worker (minion, compute) nodes.

### Control plane and worker nodes
The control plane components are the API server, the scheduler and the controller manager. 
The API server itself represents the management interface. 
The scheduler and the controller manager decide how applications should be deployed on the cluster. 
Additionally, the state and configuration of the cluster itself are controlled in the control plane components.

Worker nodes are also known as compute nodes, application nodes or minions, and are responsible for running the container workload (applications). 
The control plane for the worker nodes is implemented in the control plane components. 
The hosts running these components were historically called masters.

### Containers and images
he smallest entities in Kubernetes are Pods, which resemble your containerized application.

Using container virtualization, processes on a Linux system can be isolated up to a level where only the predefined resources are available. 
Several containers can run on the same system without “seeing” each other (files, process IDs, network). 
One container should contain one application (web server, database, cache, etc.). 
It should be at least one part of the application, e.g. when running a multi-service middleware. 
In a container itself any process can be started that runs natively on your operating system.

Containers are based on images. An image represents the file tree, which includes the binary, shared libraries and other files which are needed to run your application.

A container image is typically built from a `Containerfile` or `Dockerfile`, which is a text file filled with instructions. 
The end result is a hierarchically layered binary construct. Depending on the backend, the implementation uses overlay or copy-on-write (COW) mechanisms to represent the image.

Layer example for a Tomcat application:

1. Base image (CentOS 7)
2. Install Java
3. Install Tomcat
4. Install App

The pre-built images under version control can be saved in an image registry and can then be used by the container platform.

### Namespaces
Namespaces in Kubernetes represent a logical segregation of unique names for entities (Pods, Services, Deployments, ConfigMaps, etc.).

Permissions and roles can be bound on a per-namespace basis. This way, a user can control his own resources inside a namespace.

!!! note

    Some resources are valid cluster-wise and cannot be set and controlled on a namespace basis.

### Pods
A Pod is the smallest entity in Kubernetes.

It represents one instance of your running application process. 
The Pod consists of at least two containers, one for your application itself and another one as part of the Kubernetes design, to keep the network namespace. 
The so-called infrastructure container (or pause container) is therefore automatically added by Kubernetes.

The application ports from inside the Pod are exposed via Services.

### Services
A service represents a static endpoint for your application in the Pod. 
As a Pod and its IP address typically are considered dynamic, the IP address of the Service does not change when changing the application inside the Pod. 
If you scale up your Pods, you have an automatic internal load balancing towards all Pod IP addresses.

There are different kinds of Services:

- `ClusterIP`: Default virtual IP address range
- `NodePort`: Same as `ClusterIP` plus open ports on the nodes
- `LoadBalancer`: An external load balancer is created, only works in cloud environments, e.g. AWS ELB
- `ExternalName`: A DNS entry is created, also only works in cloud environments

A Service is unique inside a Namespace.

### Deployments
Have a look at the [official documentation](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).

### Volumes
Have a look at the [official documentation](https://kubernetes.io/docs/concepts/storage/volumes/).

### Jobs
Have a look at the [official documentation](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/).
