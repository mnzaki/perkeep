/* TODO
- what should our advertised external IP be?
  - maybe reply with IP of whatever interface the query arrived on?
	- for a setup with a bridge interface, the bridge interface IP breaks things
- reconfigure on changes to network interfaces
*/
package mdns // import "perkeep.org/internal/mdns"

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	hashiMDNS "github.com/hashicorp/mdns"
)

const (
	perkeepServiceName = "_perkeep._tcp"
	perkeepPort        = 3179
	perkeepServiceText = "Perkeep Server"
	lookupInterval     = 15 * time.Second
)

type IfaceServerMap map[string]*hashiMDNS.Server
var (
	mdnsServers IfaceServerMap
	txtRecords = []string{perkeepServiceText}
)

func StartMDNSService() func() {
	if mdnsServers == nil {
		mdnsServers = make(IfaceServerMap)
	}

	// Start a server for each interface
	// TODO FIXME XXX XXX
	// this is stupid, we should start 1 server listening on all interfaces, and
	// it should reply with different IPs depending on the origin interface of the
	// query
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		startForIface(&iface)
	}

	entriesCh := make(chan *hashiMDNS.ServiceEntry, 4)

	go func() {
		for entry := range entriesCh {
			_, isOurs := mdnsServers[entry.Addr.String()]
			if isOurs {
				continue
			}
			fmt.Printf("Got new mDNS entry: %v\n", entry)
		}
	}()
	go func() {
		for {
			hashiMDNS.Lookup(perkeepServiceName, entriesCh)
			time.Sleep(lookupInterval)
		}
	}()

	return func() {
		log.Printf("shutting down MDNS service")
		for _, server := range mdnsServers {
			server.Shutdown()
		}
		close(entriesCh)
	}
}

func startForIface(iface* net.Interface) error {
	// Setup our service export
	host, _ := os.Hostname()
	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}
	var ip net.IP
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		ip = ip.To4()
		if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
			break
		}
	}

	if ip == nil {
		return nil
	}

	service, err := hashiMDNS.NewMDNSService(host, perkeepServiceName, "", "", perkeepPort, []net.IP{ip}, txtRecords)
	var server *hashiMDNS.Server
	if err == nil {
		config := hashiMDNS.Config{Zone: service, Iface: iface}
		log.Printf("starting MDNS service for %v with IP %v", iface, ip)
		server, err = hashiMDNS.NewServer(&config)
	}

	if err != nil {
		log.Printf("Failed to start mDNS server for interface %v, err: %v", iface, err)
	} else {
		mdnsServers[ip.String()] = server
	}
	return nil
}
