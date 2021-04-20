package types

import "github.com/slack-go/slack"

type Slack struct {
	AppToken string `yaml:"app_token"`
	BotToken string `yaml:"bot_token"`
}

//incoming messageQueueElement, which will hold incoming message and metadata identifying the message.
type MessageMetadata struct {
	ChannelID string
	UserID    string
	Message   string
}

//outgoing responseQueueElement, which will hold response attachement message and metadata identifying the message
type ResponseMetadata struct {
	Attachments slack.Attachment
	Message     MessageMetadata
}

type Plugin interface {
	ProcessMessage(message *MessageMetadata) (*ResponseMetadata, error)
	Usage() string
}
