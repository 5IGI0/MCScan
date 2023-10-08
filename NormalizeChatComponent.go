package main

func NormalizeChatComponent(chat ChatComponent) string {
	text := chat.Text

	for _, subchat := range chat.Extra {
		text += NormalizeChatComponent(subchat)
	}

	return text
}
