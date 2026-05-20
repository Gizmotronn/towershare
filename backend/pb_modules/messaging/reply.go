package messaging

// BotReply is the JSON response shape for both the simulator and Twilio.
// Actions are ignored by the TwiML path.
type BotReply struct {
	Reply   string      `json:"reply"`
	Actions []BotAction `json:"actions,omitempty"`
}

type BotAction struct {
	Label string `json:"label"`
	Value string `json:"value"` // what gets sent back as the next message body
}

func reply(text string) *BotReply {
	return &BotReply{Reply: text}
}

func replyWith(text string, actions ...BotAction) *BotReply {
	return &BotReply{Reply: text, Actions: actions}
}

func act(label, value string) BotAction {
	return BotAction{Label: label, Value: value}
}
