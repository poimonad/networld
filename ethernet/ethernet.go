package ethernet

import (
	"networld/nic"
	"networld/type"
	// "networld/type"
)

type Ethernet struct {
	nic      nic.Nic
	nicE     *nic.NicE
	ip       def.Ip
	mac      def.Mac
	pendings []chan def.MacHeader
	buffer   []def.MacHeader
}

type EthernetE struct {
	nicE *nic.NicE
}

func Create(nicE *nic.NicE) EthernetE {
	return EthernetE{nicE}
}

func (ethE *EthernetE) Create(ip def.Ip, mac def.Mac) Ethernet {
	return Ethernet{ethE.nicE.Create(), ethE.nicE, ip, mac, make([]chan def.MacHeader, 0), make([]def.MacHeader, 0)}
}

func (ethE *EthernetE) Connect(eth0 *Ethernet, eth1 *Ethernet) {
	ethE.nicE.Connect(eth0.nic, eth1.nic)
}

func (eth *Ethernet) Loop() {
	go func() {
		for {
			macHeaders := <-eth.nic.Read()
			for _, macHeader := range macHeaders {
				if macHeader.IpHeader.Dst == eth.ip {
					// for arp
					eth.nic.Send(def.MacHeader{eth.mac, macHeader.Src, def.IpHeader{}})
					// fmt.Println("send arp")

				} else {
					// fmt.Println("different")
				}

				if macHeader.Dst == eth.mac {
					if len(eth.pendings) > 0 {
						eth.pendings[0] <- macHeader
						eth.pendings = eth.pendings[1:]
						// fmt.Println("resolve")
					} else {
						eth.buffer = append(eth.buffer, macHeader)
						// fmt.Println("unresoleve")
					}
				}
			}
		}
	}()
}

func (eth *Ethernet) Arp(targetIp def.Ip) def.Mac {
	go eth.nic.Send(def.MacHeader{eth.mac, def.Addr("ArpMac"), def.IpHeader{eth.ip, targetIp, ""}})
	res := <-eth.Recv()
	return res.Src
}

func (eth *Ethernet) Send(dstIp def.Ip, data string) {
	dstAddr := eth.Arp(dstIp)
	eth.nic.Send(def.MacHeader{eth.mac, dstAddr, def.IpHeader{eth.ip, dstIp, data}})
}

func (eth *Ethernet) Recv() chan def.MacHeader {
	if len(eth.buffer) > 0 {
		c := make(chan def.MacHeader, 1)
		c <- eth.buffer[0]
		eth.buffer = eth.buffer[1:]
		clear(eth.buffer)
		return c
	} else {
		c := make(chan def.MacHeader, 1)
		eth.pendings = append(eth.pendings, c)
		return c
	}
}
