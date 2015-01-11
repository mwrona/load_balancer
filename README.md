Load Balancer 
============ 

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
Run this command to download load balancer: 
``` 
go get github.com/mwrona/scalarm_load_balancer 
``` 
Now you can install it: 
```` 
go install github.com/mwrona/scalarm_load_balancer
```` 
This command will install load balancer in $GOPATH/bin. It's name will be scalarm_load_balancer 
Config 
-------- 
The config consists of config.json, cert.pem and key.pem. The cert.pem and key.pem files are needed for https server, config.json contains program configuration. 
Example of config.json:
````
{
	"MulticastAddress": "224.1.2.3:8000", 
	"CertFilePath": "cert.pem",
	"KeyFilePath": "key.pem",
	"LogDirectory" : "../log",
	"StateDirectory" : "../state",
	"RedirectionConfig" : [
		{"Path": "/", 			 "Name": "ExperimentManager"},
		{"Path": "/storage", 	 "Name": "StorageManager"},
		{"Path": "/information", "Name": "InformationService", "DisableStatusChecking": true}
	]
}

````


Description:
* LoadBalancerScheme - optional, by default: "https"; protocol on which load balancer works 
* Port - optional, by default 443 (https) or 80 (http), depends on the LoadBalancerScheme; the port on which the server listens, if port is 443 server listens also on 80 and redirects all queries to https
* MulticastAddress - address used to distribute load balancer private address
* PrivateLoadBalancerAddress - optional, by default: "localhost"; this address is send via multicast, registration to load balancer is possible only from this address and from localhost
* CertFilePath, KeyFilePath - needed only when LoadBalancerScheme is "https"; path co certificate files, by default CertFilePath: "cert.pem"; KeyFilePath: "key.pem"
* LogDirectory - optional, by default "log"; directory where logs are stored. 
* StateDirectory - optional, by default ""; directory where current state of load balancer is saved.
* RedirectionConfig - config of redirection policy: 
 * Path - path to service
 * Name - name of service type
 * DisableStatusChecking - optional, by default: false; disabling periodical status checking
 * Scheme - optional, by default: "http"; service scheme
 * StatusPath - optional, by default "/status"; path to status check
 * SecondsBetweenChecking - optional, by default 30; time beetwen periodical status checking; must be greater than zero
 * FailedConnectionsLimit - optional, by default 6; number of failed status checking before removing service; must be greater than zero

Note: MulticastAddress must be the same as in experiment manager and other services to work properly.

~~To properly work in https mode load balancer certificate must be known to all services. For development purpose you can generate self-signed certificate and install it in your system.~~

~~Instruction to generate self-signed certificate(steps 1-4): http://www.akadia.com/services/ssh_test_certificate.html~~

For now certificates checking is disable.

Run 
----

To run load balancer you should use provided in bin file scripts. For this is required that load balancer is installed in $GOPATH/bin (by go install command) and config is provided in config/config.json file. In https mode you must also provide certificate files. It is suggested to choose LogDirectory and StateDirectory same as in example_config.json.


You can also run it mannually. For that you have to supply all necessary files (config.json and in the https mode cert.pem and key.pem). By default load balancer is looking for config.json in current directory but you can specify different location as program argument. Example:
```
scalarm_load_balancer config_folder/my_config.json
```
Scalarm 
----
To run load balancer properly with Scalarm you need to run below written script (also in scalarm_registration/scalarm_registration.sh) with appropriate config. This script register load balancer in Information Service. It must be done only once.

If environment variables INFORMATION_SERVICE_URL, INFORMATION_SERVICE_LOGIN or INFORMATION_SERVICE_PASSWORD are specified they will replace config entries.

```
#!/bin/bash
#load config
source scalarm_registration_config
#script
if [ -z "$INFORMATION_SERVICE_URL" ]; then
    INFORMATION_SERVICE_URL=$INFORMATION_SERVICE_URL_CONFIG
fi
if [ -z "$INFORMATION_SERVICE_LOGIN" ]; then
    INFORMATION_SERVICE_LOGIN=$INFORMATION_SERVICE_LOGIN_CONFIG
fi
if [ -z "$INFORMATION_SERVICE_PASSWORD" ]; then
    INFORMATION_SERVICE_PASSWORD=$INFORMATION_SERVICE_PASSWORD_CONFIG
fi  
curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS" http://$INFORMATION_SERVICE_URL/experiment_managers
echo
curl -u $INFORMATION_SERVICE_LOGIN:$INFORMATION_SERVICE_PASSWORD --data "address=$REMOTE_LOAD_BALANCER_ADDRESS/storage" http://$INFORMATION_SERVICE_URL/storage_managers
echo
curl -k --data "address=$INFORMATION_SERVICE_URL&name=InformationService" https://$LOCAL_LOAD_BALANCER_ADDRESS/register
echo

```
Example of scalarm_registration_config:
```
INFORMATION_SERVICE_URL_CONFIG="localhost:11300"
INFORMATION_SERVICE_LOGIN_CONFIG="scalarm"
INFORMATION_SERVICE_PASSWORD_CONFIG="scalarm"
REMOTE_LOAD_BALANCER_ADDRESS="149.156.10.32:13585"
LOCAL_LOAD_BALANCER_ADDRESS="localhost"
```


API
-----
* /list - with parameter 'name' or without (it will print all services)
* /register - POST with parameter 'address' and 'name'
* /unregister - POST with parameter 'address' and 'name'
* /< Path > - redirection to service
