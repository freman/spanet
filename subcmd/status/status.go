package status

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/freman/spanet/pkg/spanet"
	"github.com/google/subcommands"
)

type statusCmd struct {
	spa string
}

func (*statusCmd) Name() string     { return "status" }
func (*statusCmd) Synopsis() string { return "Run the RF commamnd" }
func (*statusCmd) Usage() string {
	return `status -spa ip:port
Query the spa for it's current status and return a json blob
`
}
func (s *statusCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.spa, "spa", "", "Spa host:port")
}

func (s *statusCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	c, err := net.Dial("tcp", s.spa)
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	n := spanet.New(c)
	status, err := n.GetStatus()
	if err != nil {
		fmt.Println(err)
		return subcommands.ExitFailure
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	err = enc.Encode(status)
	if err != nil {
		panic(err)
	}

	return subcommands.ExitSuccess
}

func init() {
	subcommands.Register(&statusCmd{}, "")
}
