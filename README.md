Load Balancer 
============ 
Contents 
---------- 
* scalarm_load_balancer - main load balancer program  
* config - configuraton for load balancer 

Installation guide: 
---------------------- 
Go 
-- 
To build and install load balancer you need to install go programming language. 
You can install it from official binary distribution: 

https://golang.org/doc/install

or from source: 

https://golang.org/doc/install/source 

After that you have to specify your $GOPATH. Read more about it here: 

https://golang.org/doc/code.html#GOPATH 

Installation 
-------------- 
You can download it directly from GitHub. You have to download it into your $GOPATH/src folder 
``` 
git clone https://github.com/mwrona/scalarm_load_balancer.git 
``` 
Now run this command to download all dependencies: 
``` 
go get scalarm_load_balancer 
``` 
Now you can install load balancer: 
```` 
go install scalarm_load_balancer
```` 
This command will install load balancer in $GOPATH/bin. It's name will be scalarm_load_balancer 
Config 
-------- 
The config folder contains config.json, cert.pem and key.pem. The cert.pem and key.pem files are needed for https server, config.json contains program configuration. 
Example of config.json:
````
{
	"LocalLoadBalancerAddress": "localhost:9000",
	"RemoteBalancerAddress": "localhost:9000",
	"Port": "9000",
	"MulticastAddress": "224.1.2.3:8000", 
	"LoadBalancerScheme": "https",
	"CertificateCheckDisable": true,
	"InformationServiceAddress": "localhost:11300",
	"InformationServiceScheme": "http",
	"InformationServiceUser" : "scalarm",
	"InformationServicePass" : "scalarm",
	"LoadBalancerUser" : "scalarm",
	"LoadBalancerPass" : "scalarm",
	"CertFilePath": "cert.pem",
	"KeyFilePath": "key.pem"
}

````
Note: MulticastAddress must be the same as in experiment manager and other services to work properly.

Run 
---- 
To run load balancer you have to supply all necessary files (config.json and in the https mode cert.pem and key.pem). By default load balancer is looking for config.json in current directory but you can specify different location as program argument. Example:
```
scalarm_load_balancer config_folder/my_config.json
```
API
-----
* / - redirection to Experimet Managers
* /information - redirection to Information Services
* /storage - redirection to Storage Managers
* /experiment_managers - list of available Experiments Managers
* /experiment_managers/register - POST with parameter address, registration of new Experiment Manager
* /experiment_managers/unregister - POST with parameter address, unregistration of Experiment Manager
* /storage_managers - list of available Storage Managers
* /storage_managers/register - POST with parameter address, registration of new Storage Managers
* /storage_managers/unregister - POST with parameter address, unregistration of Storage Managers
