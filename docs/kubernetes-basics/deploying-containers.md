# Deploying Containers
In this tutorial, we are going to deploy our first container image and look at the concepts of Pods, Services, and Deployments.

## Task: Start and stop a single Pod
After we’ve familiarized ourselves with the platform, we are going to have a look at deploying a pre-built container image or any other public container registry.

First, we are going to directly start a new Pod.
For this we have to define our Kubernetes Pod resource definition. 
Create a new file `03_pod.yaml` with the following content: