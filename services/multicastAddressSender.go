package services

import (
	"net"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/utils"
	"time"

	"code.google.com/p/go.net/ipv4"
)

func MulticastAddressSender(proxyAddress, proxyPort string) {
	config, err := model.LoadConfig()
	utils.Check(err)

	mcaddr, err := net.ResolveUDPAddr("udp", config.Address)
	utils.Check(err)

	// conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	// utils.Check(err)

	c, err := net.ListenPacket("udp4", "")
	utils.Check(err)
	defer c.Close()

	conn := ipv4.NewPacketConn(c)

	conn.JoinGroup(nil, mcaddr)
	conn.SetMulticastLoopback(true)

	b := make([]byte, 20)
	copy(b, proxyAddress+":"+proxyPort)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			// _, err = conn.WriteToUDP(b, mcaddr)
			_, err = conn.WriteTo(b, nil, mcaddr)
			utils.Check(err)
		}
	}
}
