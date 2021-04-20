package types

type Slack struct {
	AppToken string `yaml:"app_token"`
	BotToken string `yaml:"bot_token"`
}

type MessageMetadata struct {
	ChannelID string
	UserID    string
	Message   string
	Timestamp string
}

type Plugin interface {
	ProcessMessage(b Backend, message *MessageMetadata) (*MessageMetadata, error)
	Usage() string
}
