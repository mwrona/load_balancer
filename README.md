Load Balancer 
============ 
Contents 
---------- 
* reverseProxy - main load balancer program 
* server - experiment manager mock 
* client - client mock,  
* config - configuraton for load balancer 

Note: server and client may not work with current version of load balancer. They may require some adjustments. 

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
You can download it directly from GitHub. You have to downlaod it into your $GOPATH/src folder 
``` 
git clone https://github.com/mwrona/load_balancer.git 
``` 
Now run this command to download all dependencies: 
``` 
go get load_balancer/reverseProxy 
``` 
Now you can install load balancer: 
```` 
go install load_balancer/reverseProxy 
```` 
This command will install load balancer in $GOPATH/bin. It's name will be reverseProxy 
Build Options 
---------------- 
With -tags option you can specify build options:  
* no parameter: http server 
* prod : https server 
* certOff: disabling certificate checking for https 

Paramters can be mixed. For example: 
``` 
go install -tags "prod certOff" load_balancer/reverseProxy 
``` 
Note: Use -a option in go install if you didn't change any files after previous install. 
Config 
-------- 
The config folder contains config.txt, cert.pem and key.pem. Cert.pem and key.pem are needed for https server, config.txt contains program configuration. For now it's only multicast address. We recommend not to change that, the same address is used in experiment manager. 
Run 
---- 
To run you have to copy contents of config folder to folder with executable of load balancer. By default it will be $GOPATH/bin 
