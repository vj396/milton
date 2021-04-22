package slack

import (
	"fmt"
	"strings"
	"sync"

	slacklib "github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/vj396/milton/src/types"
)

var (
	registry  = make(map[string]types.Plugin)
	lock      = new(sync.RWMutex)
	apiClient = new(slacklib.Client)
)

func Client(slackConfig *types.Slack) (*socketmode.Client, error) {
	if slackConfig.AppToken == "" {
		return nil, fmt.Errorf("app token is empty")
	}
	if !strings.HasPrefix(slackConfig.AppToken, "xapp-") {
		return nil, fmt.Errorf("app_token must have the prefix \"xapp-\"")
	}
	if slackConfig.BotToken == "" {
		return nil, fmt.Errorf("bot token is empty")
	}
	if !strings.HasPrefix(slackConfig.BotToken, "xoxb-") {
		return nil, fmt.Errorf("bot_token must have the prefix \"xoxb-\"")
	}
	apiClient = slacklib.New(
		slackConfig.BotToken,
		//slack.OptionDebug(true),
		//slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slacklib.OptionAppLevelToken(slackConfig.AppToken),
	)
	socketmodeClient := socketmode.New(apiClient)

	return socketmodeClient, nil
}

func Register(command string, plugin types.Plugin) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := registry[command]; ok {
		return fmt.Errorf("command: %+q already registered", command)
	}
	registry[command] = plugin
	return nil
}

func GetRegistry() map[string]types.Plugin {
	return registry
}

func ApiClient() *slacklib.Client {
	return apiClient
}
