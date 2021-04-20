package opsgenie

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type opsgenieAddCmd struct{}

func init() {
	o := new(opsgenieAddCmd)
	err := slack.Register("!opsgenie-add", o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (o *opsgenieAddCmd) ProcessMessage(b types.Backend, message *types.MessageMetadata) *types.MessageMetadata {
	return nil
}

func (o *opsgenieAddCmd) Usage() string {
	return "<!opsgenie-add 111-222-333>"
}
