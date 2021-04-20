package interrupts

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type doneCmd struct{}

func init() {
	d := new(doneCmd)
	err := slack.Register("!done", d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (dc *doneCmd) ProcessMessage(message *types.MessageMetadata) (*types.ResponseMetadata, error) {
	return nil, nil
}

func (dc *doneCmd) Usage() string {
	return `Usage <!done D123456> || <!done D123456,D234567,https://example.com/D345678>`
}
