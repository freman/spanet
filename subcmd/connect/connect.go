package connect

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/freman/spanet/pkg/flags"
	"github.com/freman/spanet/pkg/wifly"
	"github.com/google/subcommands"
)

var spanet = &net.IPNet{IP: net.IP{1, 2, 3, 0}, Mask: net.IPMask{255, 255, 255, 0}}

type connectCmd struct {
	targetIP flags.IP
	ssid     string
	password string
}

func (*connectCmd) Name() string     { return "connect" }
func (*connectCmd) Synopsis() string { return "Connect the spa to your network" }
func (*connectCmd) Usage() string {
	return `connect [-target targetip] -ssid {ssid} -password {password}:
	Connect the spa to your network, using -target ip skips the initial wifi connect step

`
}
func (c *connectCmd) SetFlags(f *flag.FlagSet) {
	c.targetIP.Default = net.IP{1, 2, 3, 4}

	f.Var(&c.targetIP, "target", "Target IP")
	f.StringVar(&c.ssid, "ssid", "", "SSID to connect to")
	f.StringVar(&c.password, "password", "", "Password to connect with")
}

func (c *connectCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.ssid == "" || c.password == "" {
		fmt.Println("ssid and password are required")
		return subcommands.ExitUsageError
	}

	if !c.targetIP.IsSet() {
		waitForNetwork()
	}

	fmt.Println("Connecting to Spalink")
	conn, err := net.Dial("tcp", c.targetIP.IP().String()+":2000")
	if err != nil {
		fmt.Println("Failed to connect to spanet controller")
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	op, err := wifly.NewCommand(conn)
	if err != nil {
		fmt.Println("Failed to enter command mode")
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println("Scanning for access points")
	// Could just use "join" but I'd rather not have to reset the damn adapter
	aps, err := op.Scan()
	if err != nil {
		fmt.Println("Error while scanning for access points")
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	var matched []wifly.WIFIRecord
	for _, v := range aps {
		if v.SSID == c.ssid {
			matched = append(matched, v)
		}
	}

	if len(matched) == 0 {
		fmt.Printf("Unable to find ssid %q\n", c.ssid)
		return subcommands.ExitFailure
	}

	// A lot of effort to go through but if someone has multiple AP"s we"ll find the loudest one
	sort.Slice(matched, func(i, j int) bool { return matched[i].RSSI > matched[j].RSSI })

	replacement, ssid, password := wifly.ReplaceSpaces(c.ssid, c.password)

	fmt.Println("Configuring Spalink")

	op.Set("ip dhcp", "1")

	if replacement != wifly.DefaultReplacementCharacter {
		op.Set("opt replace", fmt.Sprintf("0x%x", replacement[0]))
	}

	op.Set("wlan ssid", ssid)
	op.Set("wlan phrase", password)

	if replacement != wifly.DefaultReplacementCharacter {
		op.Set("opt replace", fmt.Sprintf("0x%x", wifly.DefaultReplacementCharacter[0]))
	}

	op.Set("wlan channel", "0")
	op.Set("wlan auth", strconv.Itoa(matched[0].Auth))
	op.Set("wlan linkmon", "30")
	op.Set("comm idle", "300")
	op.Set("comm remote", "0")

	// Disable autoreconn and ip host, should leave it listening on port 2000 on the lan without trying to phone anyone
	op.Set("sys autoconn", "0")
	op.Set("ip host", "0.0.0.0")
	op.Set("wlan join", "1")

	if err := op.Save(); err != nil {
		fmt.Println("Error saving config")
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	fmt.Println("Done, rebooting")
	op.Reboot()

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
