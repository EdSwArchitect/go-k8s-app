# Another Kubernetes Play Application

That's all this is. **Playing around**.

## Ingester

An app that reads files and sends them to ElasticSearch for indexing. It makes use of the k8s
configmap in order to reference the configuration file, which contains information about the
rest of playground application.

The Docker file **Docker** is used to create the Docker image edswarchitect/k8-ingester.

To build the app:

```
cd ingester
go build
cd ..
```

To build the Docker image.

```
docker build . -t edswarchitect/k8s-ingester
docker push edswarchitect/k8s-ingester
```

### Associated k8s files

***ingester-configmap-pod.yaml*** - Creates the pod for the ingester app.

***ingester-configmap-svc.yaml*** - Creates the service for the ingester app.


## Queryer
Another app that retrieves documents from ElasticSearch via REST calls.


## Default Application

All this does is handle "unregistered" URIs.

The Docker file **DockerDefault** is used to build it.

To build the app.

```
cd default-app
go build
cd ..
```

To build the Docker image.

```
docker build . -t edswarchitect/default-app -f DockerfileDefault
docker push edswarchitect/default-app
```

### Associated k8s files

***ingester-defaultapp-pod.yaml*** - Creates the pod for the default web app.

***ingester-defaultapp-svc.yaml*** - Creates the service for the default web app.


## Kubernetes

This is run in the Google Cloud Kubernetes

```
gcloud container clusters get-credentials edwin --zone us-east1-b --project learn-vm
```

From the command line to connect to the Kubernetes cluster. The cluster name
was edwin, in the zone us-east1-b, in the project lean-vm


### Directory config

Contains the information that will be place into a configmap that is expected
to be named **ingester-configmap**. How the configmap is created.

```
kubectl create configmap ingester-config  --from-file=config
```


## ElasticSearch

This sets up the ElasticSearch pod and service.


### Associated k8s files

***es-pod.yaml*** - Creates the pod for ElasticSearch.

***ingester-es-svc.yaml*** - Creates the service for ElasticSearch.
