package support

import (
	"net"
	"os"
)

func GetMulticastGroupAddrFromEnv() *net.UDPAddr {
	addr, okAddr := os.LookupEnv("MULTICAST_ADDR")
	port, okPort := os.LookupEnv("MULTICAST_PORT")

	if okAddr && okPort {
		udpAddr, err := net.ResolveUDPAddr("udp", addr + ":" + port)
		if err != nil {
			panic(err)
		}
		return udpAddr
	}
	return nil
}
