package connect

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"git.fremnet.net/spanet/pkg/flags"
	"github.com/google/subcommands"
)

var spanet = &net.IPNet{IP: net.IP{1, 2, 3, 0}, Mask: net.IPMask{255, 255, 255, 0}}

type connectCmd struct {
	targetIP flags.IP
}

func (*connectCmd) Name() string     { return "connect" }
func (*connectCmd) Synopsis() string { return "Connect the spa to your network" }
func (*connectCmd) Usage() string {
	return `connect [-target targetip]:
	Connect the spa to your network, using -target ip skips the initial wifi connect step

`
}
func (c *connectCmd) SetFlags(f *flag.FlagSet) {
	c.targetIP.Default = net.IP{1, 2, 3, 4}

	f.Var(&c.targetIP, "target", "Target IP")
}

func (c *connectCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if !c.targetIP.IsSet() {
		waitForNetwork()
	}

	return subcommands.ExitSuccess
}

func waitForNetwork() {
	fmt.Println("Set your spa WIFI to 'HOT' mode")
	fmt.Println("Connect to your spa's wifi")

	for {
		fmt.Print(".")
		ifaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		// handle err
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				panic(err)
			}
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// process IP address
				if spanet.Contains(ip) {
					fmt.Println("\nspa net detected")
					return
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func init() {
	subcommands.Register(&connectCmd{}, "")
}
