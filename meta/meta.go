package meta

import (
	"fmt"
	def "networld/type"
)

type Connection struct {
	Src int
	Dst int
}

type RawDevice struct {
	Id     int
	Buffer []def.MacHeader
}

func (dev *RawDevice) write(data def.MacHeader) {
	dev.Buffer = append(dev.Buffer, data)
}

type Await struct {
	id int
	ch chan []def.MacHeader
}

type Meta struct {
	connections []Connection
	devices     []RawDevice
	nextId      int
	awaits      []Await
}

func Create() Meta {
	return Meta{
		connections: []Connection{},
		devices:     []RawDevice{},
		nextId:      0,
		awaits:      []Await{},
	}
}

func (meta *Meta) Create() int {
	new_id := meta.nextId
	meta.nextId++
	meta.devices = append(meta.devices, RawDevice{new_id, make([]def.MacHeader, 0)})
	return new_id
}

func (meta *Meta) Device(id int) (*RawDevice, error) {
	for i, device := range meta.devices {
		if device.Id == id {
			return &meta.devices[i], nil
		}
	}
	return nil, fmt.Errorf("device %d not found", id)
}

func (meta *Meta) Connect(src, dst int) {
	meta.connections = append(meta.connections, Connection{src, dst})
}

func (meta *Meta) Broadcast(src int, data def.MacHeader) {
	for _, connection := range meta.connections {
		if connection.Src == src {
			delivered := false
			new_awaits := make([]Await, 0, len(meta.awaits))

			for _, await := range meta.awaits {
				if !delivered && await.id == connection.Dst {
					await.ch <- []def.MacHeader{data}
					delivered = true
					continue
				}
				new_awaits = append(new_awaits, await)
			}

			meta.awaits = new_awaits

			if !delivered {
				if dev, err := meta.Device(connection.Dst); err == nil {
					dev.write(data)
				}
			}
		}
	}
}

func (meta *Meta) Receive(src int, c chan []def.MacHeader) {
	dev, err := meta.Device(src)
	if err != nil {
		return
	}

	if len(dev.Buffer) > 0 {
		received := make([]def.MacHeader, len(dev.Buffer))
		copy(received, dev.Buffer)
		clear(dev.Buffer)
		c <- received
	} else {
		meta.awaits = append(meta.awaits, Await{id: src, ch: c})
	}
}
