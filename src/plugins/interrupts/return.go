package interrupts

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type returnCmd struct{}

func init() {
	r := new(returnCmd)
	err := slack.Register("!return", r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (rc *returnCmd) ProcessMessage(message *types.MessageMetadata) (*types.ResponseMetadata, error) {
	return nil, nil
}

func (rc *returnCmd) Usage() string {
	return `Usage <!return D123456> || <!return D123456,D234567,https://example.com/D345678>`
}
