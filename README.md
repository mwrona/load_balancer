Load Balancer 
============ 
Contents 
---------- 
* scalarm_load_balancer - main load balancer program  
* config - example configuraton for load balancer 

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
The config consists of config.json, cert.pem and key.pem. The cert.pem and key.pem files are needed for https server, config.json contains program configuration. 
Example of config.json:
````
{
	"LocalLoadBalancerAddress": "localhost",
	"RemoteLoadBalancerAddress": "149.156.10.32:13585",
	"Port": "443",
	"MulticastAddress": "224.1.2.3:8000", 
	"LoadBalancerScheme": "https",
	"InformationServiceAddress": "localhost:11300",
	"InformationServiceUser" : "scalarm",
	"InformationServicePass" : "scalarm",
	"CertFilePath": "cert.pem",
	"KeyFilePath": "key.pem",
	"RedirectionConfig" : [
		{"Path": "/", 			 "Name": "ExperimentManager"},
		{"Path": "/storage", 	 "Name": "StorageManager"},
		{"Path": "/information", "Name": "InformationService", "DisableStatusChecking": true, "Scheme": "http"}
	]
}

````
Note: MulticastAddress must be the same as in experiment manager and other services to work properly.

Optional entries:
* LocalLoadBalancerAddress - default: "localhost"
* LoadBalancerScheme - defaulf: "https"
* CertFilePath, KeyFilePath when LoadBalancerScheme is "http"
* In RedirectionConfig: 
 * DisableStatusChecking - default: false
 * Scheme - default: "http"

If environment variables INFORMATION_SERVICE_URL, INFORMATION_SERVICE_LOGIN or INFORMATION_SERVICE_PASSWORD are specified they will replace config entries. In this case config entries (InformationServiceAddress, InformationServiceUser, InformationServicePass) can be omitted.

Run 
---- 
To run load balancer you have to supply all necessary files (config.json and in the https mode cert.pem and key.pem). By default load balancer is looking for config.json in current directory but you can specify different location as program argument. Example:
```
scalarm_load_balancer config_folder/my_config.json
```
API
-----
* /list - with parameter 'name'
* /register - POST with parameter 'address' and 'name'
