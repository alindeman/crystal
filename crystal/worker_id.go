package crystal

import (
	"net"
)

func WorkerIdFromNetworkInterfaceName(name string) (id WorkerId, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return id, err
	}

	copy(id[:], iface.HardwareAddr[0:len(id)])
	return id, nil
}
