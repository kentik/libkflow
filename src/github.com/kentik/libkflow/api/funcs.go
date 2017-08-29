package api

import (
	"net"
	"strings"
)

func NormalizeName(name string) string {
	return strings.Replace(name, ".", "_", -1)
}

func ActiveNetworkInterfaces() (map[string]InterfaceUpdate, error) {
	active := map[string]InterfaceUpdate{}

	nifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, nif := range nifs {
		up := nif.Flags&net.FlagUp != 0

		addrs, err := interfaceAddrs(&nif)
		if err != nil {
			continue
		}

		if !up || len(addrs) == 0 || nif.HardwareAddr == nil {
			continue
		}

		addr := addrs[0]
		if len(addrs) > 1 {
			addrs = addrs[1:]
		} else {
			addrs = nil
		}

		n := len(nif.HardwareAddr)
		a := int(nif.HardwareAddr[n-2])
		b := int(nif.HardwareAddr[n-1])

		active[nif.Name] = InterfaceUpdate{
			Index:   uint64(a<<8 | b),
			Alias:   "",
			Desc:    nif.Name,
			Address: addr.Address,
			Netmask: addr.Netmask,
			Addrs:   addrs,
		}
	}

	return active, nil
}

func interfaceAddrs(nif *net.Interface) ([]Addr, error) {
	all, err := nif.Addrs()
	if err != nil {
		return nil, err
	}

	var addrs []Addr
	for _, a := range all {
		if a, ok := a.(*net.IPNet); ok && a.IP.IsGlobalUnicast() {
			addrs = append(addrs, Addr{
				Address: a.IP.String(),
				Netmask: net.IP(a.Mask).String(),
			})
		}
	}

	return addrs, nil
}
