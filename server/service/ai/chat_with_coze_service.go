package ai

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/spf13/cast"
)

type ChatWithCozeService struct {
}

var ChatWithCozeServiceApp = &ChatWithCozeService{}

func (cwc *ChatWithCozeService) ChatWithCoze(message string, history []*CozeChatMessage) (string, []*CozeChatMessage, error) {
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + global.GVA_CONFIG.Coze.AccessToken
	//header["Content-Type"] = "application/json"
	header["Connection"] = "keep-alive"
	header["Accept"] = "*/*"
	header["Host"] = "api.coze.cn"

	apiUrl := global.GVA_CONFIG.Coze.ApiUrl
	//refreshToken := chatModel.RefreshToken

	chatReq := &CozeChatReq{
		BotId:       global.GVA_CONFIG.Coze.BotId,
		User:        cast.ToString(utils.GenID()),
		Query:       message,
		Stream:      false,
		ChatHistory: history,
	}

	history = append(history, &CozeChatMessage{
		Role:        "user",
		Content:     message,
		ContentType: "text",
	})

	statusCode, respBody, err := utils.PostWithHeaders(apiUrl, utils.Parse2Json(chatReq), header)

	if err != nil {
		return "", history, err
	}

	if statusCode != 200 {
		return "", history, errors.New(string(respBody))
	}

	chatResp := &CozeChatResp{}
	err = utils.JsonStrToStruct(string(respBody), chatResp)
	if err != nil {
		return "", history, err
	}

	if len(chatResp.Messages) > 0 {
		history = append(history, &CozeChatMessage{
			Role:        "assistant",
			Content:     chatResp.Messages[0].Content,
			Type:        "answer",
			ContentType: "text",
		})
		return chatResp.Messages[0].Content, history, err
	}

	return "", history, errors.New("chat回复为空")
}

type CozeChatReq struct {
	BotId          string             `json:"bot_id"`
	ConversationId string             `json:"conversation_id"`
	User           string             `json:"user"`
	Query          string             `json:"query"`
	ChatHistory    []*CozeChatMessage `json:"chat_history,omitempty"`
	Stream         bool               `json:"stream"`
}

type CozeChatResp struct {
	Messages       []CozeChatMessage `json:"messages"`
	ConversationId string            `json:"conversation_id"`
	Code           int               `json:"code"`
	Msg            string            `json:"msg"`
}

type CozeChatMessage struct {
	Role        string `json:"role"`
	Type        string `json:"type,omitempty"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	ExtraInfo   string `json:"extra_info,omitempty"`
}
