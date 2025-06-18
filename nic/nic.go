package nic

import (
	"networld/meta"
	def "networld/type"
)

type Nic struct {
	device int
	nicE   *NicE
}

type NicE struct {
	meta *meta.Meta
}

func Create(meta *meta.Meta) NicE {
	return NicE{meta}
}

func (nicE *NicE) Create() Nic {
	raw_device := nicE.meta.Create()
	return Nic{raw_device, nicE}
}

func (nicE *NicE) Connect(n0 Nic, n1 Nic) {
	nicE.meta.Connect(n0.device, n1.device)
}

func (nic *Nic) Send(data def.MacHeader) {
	go nic.nicE.meta.Broadcast(nic.device, data)
}

func (nic *Nic) Read() chan []def.MacHeader {
	c := make(chan []def.MacHeader)
	go nic.nicE.meta.Receive(nic.device, c)
	return c
}
