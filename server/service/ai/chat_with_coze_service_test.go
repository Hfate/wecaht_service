package ai

import (
	"testing"
)

func TestChatWithCozeService_ChatWithCoze(t *testing.T) {
	message := "河南两家医院涉欺诈骗保被查处"
	history := make([]*CozeChatMessage, 0)
	ChatWithCozeServiceApp.ChatWithCoze(message, history)
}
