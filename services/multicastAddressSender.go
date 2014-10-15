package services

import (
	"log"
	"net"
	"scalarm_load_balancer/utils"
	"time"

	"code.google.com/p/go.net/ipv4"
)

func StartMulticastAddressSender(loadBalancerAddress, multicastAddress string) {
	c := make(chan error)

	for {
		if _, err := utils.RepetitiveCaller(
			func() (interface{}, error) {
				go multicastAddressSender(loadBalancerAddress, multicastAddress, c)
				err := <-c
				return nil, err
			},
			[]int{5, 5, 10, 10, 30},
			"MulticastAddressSender",
		); err != nil {
			log.Fatal("Unable to send address via multicast, stopping load balancer")
		}

		err := <-c
		log.Printf("MulticastAddressSender: an error occured:\n%s\nTrying to restart", err.Error())
	}
}

func multicastAddressSender(loadBalancerAddress, multicastAddress string, out chan error) {
	mcaddr, err := net.ResolveUDPAddr("udp", multicastAddress)
	if err != nil {
		out <- err
		return
	}
	// conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	// utils.Check(err)

	c, err := net.ListenPacket("udp4", "")
	if err != nil {
		out <- err
		return
	}
	defer c.Close()

	conn := ipv4.NewPacketConn(c)

	err = conn.JoinGroup(nil, mcaddr)
	if err != nil {
		out <- err
		return
	}

	err = conn.SetMulticastLoopback(true)
	if err != nil {
		out <- err
		return
	}

	b := make([]byte, 20)
	copy(b, loadBalancerAddress)

	_, err = conn.WriteTo(b, nil, mcaddr)
	if err != nil {
		out <- err
		return
	}

	out <- nil

	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			// _, err = conn.WriteToUDP(b, mcaddr)
			_, err = conn.WriteTo(b, nil, mcaddr)
			if err != nil {
				out <- err
				return
			}
		}
	}
}
