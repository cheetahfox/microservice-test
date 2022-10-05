# microservice-test
A small example microservice test designed to run In a K8s Cluster.

# Building the docker image. 
This image should be easy to build. 

## Requirements 
It should work with most recent docker versions but tested working with "Docker version 20.10.12". 

You must be able to build docker images. This will require sudo permissions or running the build as root.

## Docker Image Build commands 
First clone from Github. 

````
git clone https://github.com/cheetahfox/microservice-test.git
````

Then run the following docker build command to build the image locally.

````
sudo docker build --no-cache -t cheetahfox/microservice-project:0.02 .
````

To run it just from docker use the following command while in the cloned repo and after you have built it locally. 
````
sudo docker run --publish 2200:2200 --env-file TestENV.sh cheetahfox/microservice-project:0.02
````

# Running in Kubernetes
I have included a deployment.yaml file that creates the following Kubernetes resources. Operation in Kubernetes is the perfered method of running this project. 

### Resources
* configmap/microservice-config
* secret/ms-secret  
* deployment/microservice
* service/microservice
* ingress/microservice

You will need to customize the configuration options in the deployment.yaml file as the examples included have specific settings uniquie to my local setup. Specifically, the following items may need to be updated with values customized for your deployment.

### Settings to customize
* configmap : SYMBOL - This is the symbol to query.
* configmap : NDAYS  - This is the number of days to query (must be between 1-100).
* secret: API_KEY - This is the API Key for AlphaVantage and should be changed (see for more info. https://www.alphavantage.co/support/#support) 
* ingress: ingressClassName - Currently I have ingress-nginx in my enviroment, you may need to adjust this if your env is different.
* ingress: host - This is the DNS host for the service. It must resovle to the ip address of your ingress controller (see ingressClassName above)
