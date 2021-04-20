package interrupts

import (
	"fmt"
	"os"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type queueCmd struct{}

func init() {
	q := new(queueCmd)
	err := slack.Register("!queue", q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (qc *queueCmd) ProcessMessage(message *types.MessageMetadata) (*types.ResponseMetadata, error) {
	return nil, nil
}

func (qc *queueCmd) Usage() string {
	return `Usage <!queue D123456> || <!queue D123456,D234567,https://example.com/D345678>`
}
