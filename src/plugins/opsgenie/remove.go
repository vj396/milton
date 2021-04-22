package opsgenie

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type opsgenieRmCmd struct{}

func init() {
	o := new(opsgenieRmCmd)
	err := slack.Register("!opsgenie", o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (o *opsgenieRmCmd) ProcessMessage(b types.Backend, message *types.MessageMetadata) *types.MessageMetadata {
	return nil
}

func (o *opsgenieRmCmd) Usage() string {
	return "<!opsgenie-rm 111-222-333>"
}
