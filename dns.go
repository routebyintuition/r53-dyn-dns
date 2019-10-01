package main

import (
	"io/ioutil"
	"net"
	"net/http"
)

// GetPublicIPService returns the public ip address using external service
func (conf *Config) GetPublicIPService() (net.IP, error) {
	resp, err := http.Get(conf.DNSURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(ip)), nil
}

// CheckIPMatching does a DNS lookup on the configured hostname to check the publicly resolving IP of the hostname
// reduce the Route53 DNS TTL to get faster updates
// Returns error if there is an error
// Returns TRUE if the
func (conf *Config) CheckIPMatching(serviceIP *net.IP, dnsIP *net.IP) (bool, error) {
	ips, err := net.LookupIP(conf.Hostname)
	if err != nil {
		return false, err
	}

	for _, ip := range ips {
		Info.Printf("%s -> %s\n", conf.Hostname, ip.String())
	}

	return true, nil
}

// GetDNSEntries returns all of the IPs associated with the hostname
func GetDNSEntries(conf Config) ([]net.IP, error) {
	ips, err := net.LookupIP(conf.Hostname)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

// PrintIPList Pretty pring of current DNS entries for STDOUT and logs
func PrintIPList(conf Config) {
	ips, err := GetDNSEntries(conf)
	if err != nil {
		Error.Println("Could not get DNS information...will try again on next loop...")
		return
	}

	for _, ip := range ips {
		Info.Printf("%s -> %s\n", conf.Hostname, ip.String())
	}
}
