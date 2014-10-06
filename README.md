Load Balancer 
============ 
Contents 
---------- 
* scalarm_load_balancer - main load balancer program  
* config - example configuraton for load balancer 
* scripts - scripts to start on stop load balancer (on linux)

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
	"PublicLoadBalancerAddress": "149.156.10.32:13585",
	"Port": "443",
	"MulticastAddress": "224.1.2.3:8000", 
	"CertFilePath": "cert.pem",
	"KeyFilePath": "key.pem",
	"RedirectionConfig" : [
		{"Path": "/", 			 "Name": "ExperimentManager"},
		{"Path": "/storage", 	 "Name": "StorageManager"},
		{"Path": "/information", "Name": "InformationService", "DisableStatusChecking": true, "Scheme": "http"}
	]
}

````


Description:
* PrivateLoadBalancerAddress - optional, by default: "localhost"; this address is send via multicast
* PublicLoadBalancerAddress - this address is registered in Information Service
* Port - the port on which the server listens, if port is 443 server listens also on 80 and redirects all queries to https
* MulticastAddress - address used to distribute load balancer private address
* LoadBalancerScheme - optional, by default: "https"; protocol on which load balancer works 
* CertFilePath, KeyFilePath - needed only when LoadBalancerScheme is "https"; path co certificate files, by default CertFilePath: "cert.pem"; KeyFilePath: "key.pem"
* RedirectionConfig - config of redirection policy: 
 * Path - path to service
 * Name - name of service type
 * DisableStatusChecking - optional, by default: false; disabling periodical status checking
 * Scheme - optional, by default: "http"; service scheme

Note: MulticastAddress must be the same as in experiment manager and other services to work properly.

If environment variables INFORMATION_SERVICE_URL, INFORMATION_SERVICE_LOGIN or INFORMATION_SERVICE_PASSWORD are specified they will replace config entries. In this case config entries (InformationServiceAddress, InformationServiceUser, InformationServicePass) can be omitted.

~~To properly work in https mode load balancer certificate must be known to all services. For development purpose you can generate self-signed certificate and install it in your system.~~

~~Instruction to generate self-signed certificate(steps 1-4): http://www.akadia.com/services/ssh_test_certificate.html~~

For now certificates checking is disable.

Run 
----
To run load balancer you have to supply all necessary files (config.json and in the https mode cert.pem and key.pem). By default load balancer is looking for config.json in current directory but you can specify different location as program argument. Example:
```
scalarm_load_balancer config_folder/my_config.json
```
To run properly with Scalarm you need to run below written script (scrpits/scalarm_registration.sh) with appropriate config after first run. You have to run it only once.

```
#!/bin/bash
#config
INFORMATION_SERVICE_URL_="localhost:11300"
INFORMATION_SERVICE_LOGIN_="scalarm"
INFORMATION_SERVICE_PASSWORD_="scalarm"
REMOTE_LOAD_BALANCER_ADDRESS="149.156.10.32:13585"
LOCAL_LOAD_BALANCER_ADDRESS="localhost"
#script
if [ -z "$INFORMATION_SERVICE_URL" ]; then
    INFORMATION_SERVICE_URL=$INFORMATION_SERVICE_URL_
fi
if [ -z "$INFORMATION_SERVICE_LOGIN" ]; then
    INFORMATION_SERVICE_LOGIN=$INFORMATION_SERVICE_LOGIN_
fi
if [ -z "$INFORMATION_SERVICE_PASSWORD" ]; then
    INFORMATION_SERVICE_PASSWORD=$INFORMATION_SERVICE_PASSWORD_
fi  
curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS" http://$INFORMATION_SERVICE_URL/experiment_managers
echo
curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS/storage" http://$INFORMATION_SERVICE_URL/storage_managers
echo
curl -k --data "address=$INFORMATION_SERVICE_URL&name=InformationService" https://$LOCAL_LOAD_BALANCER_ADDRESS/register
echo

```


API
-----
* /list - with parameter 'name'
* /register - POST with parameter 'address' and 'name'
