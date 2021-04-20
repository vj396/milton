package opsgenie

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type opsgenieSearchCmd struct{}

func init() {
	o := new(opsgenieSearchCmd)
	err := slack.Register("!opsgenie", o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (o *opsgenieSearchCmd) ProcessMessage(b types.Backend, message *types.MessageMetadata) *types.MessageMetadata {
	return nil
}

func (o *opsgenieSearchCmd) Usage() string {
	return "<!opsgenie SRE> || <!opsgenie infra>"
}
