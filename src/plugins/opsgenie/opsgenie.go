package opsgenie

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type opsgenieCmd struct{}

func init() {
	o := new(opsgenieCmd)
	err := slack.Register("!opsgenie", o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (o *opsgenieCmd) ProcessMessage(message *types.MessageMetadata) (*types.ResponseMetadata, error) {
	return nil, nil
}

func (o *opsgenieCmd) Usage() string {
	return `Usage <!opsgenie add 111-222-333> || <!opsgenie remove 111-222-333> || <!opsgenie list> || <!opsgenie SRE> || <!opsgenie infra>`
}
