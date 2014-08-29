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
This command will install load balancer in $GOPATH/bin. It's name will be reverseProxy 
Config 
-------- 
The config folder contains config.json, cert.pem and key.pem. The cert.pem and key.pem files are needed for https server, config.json contains program configuration. 
Example of config.json:
````
{
	"LoadBalancerAddress": "localhost:9000",
	"Port": "9000",
	"MulticastAddress": "224.1.2.3:8000", 
	"LoadBalancerScheme": "https",
	"CertificateCheckDisable": true,
	"InformationServiceAddress": "localhost:11300",
	"InformationServiceScheme": "http",
	"CertFilePath": "cert.pem",
	"KeyFilePath": "key.pem"
}

````
Note: MulticastAddress must be the same as in experiment manager and other services to work properly.

Run 
---- 
To run you have to copy contents of config folder to folder with executable of load balancer. By default it will be $GOPATH/bin 
