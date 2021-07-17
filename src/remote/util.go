package main

import (
	"errors"
	"net"
)

func findNetworkDevices() []net.IP {
	ifaces, err := net.Interfaces()

	if err != nil {
		return devices
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			devices = append(devices, ip)
		}
	}

	return devices
}

func findPackageById(id string) (Package, error) {
	for _, item := range packages {
		if item.ID == id {
			return item, nil
		}
	}

	return Package{}, errors.New("file does not exist")
}
