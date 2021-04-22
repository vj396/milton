package help

import (
	"fmt"
	"os"
	"strings"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type helpCmd struct{}

func init() {
	h := new(helpCmd)
	err := slack.Register("!help", h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (hc *helpCmd) ProcessMessage(b types.Backend, message *types.MessageMetadata) *types.MessageMetadata {
	registry := slack.GetRegistry()
	var payload []string
	for cmd, plugin := range registry {
		r := fmt.Sprintf("%s\t-\t%s", cmd, plugin.Usage())
		payload = append(payload, r)
	}
	message.Message = strings.Join(payload, "\n")
	return message
}

func (hc *helpCmd) Usage() string {
	return "<!help>"
}
