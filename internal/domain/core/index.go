package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type St struct {
	discordWebhookUrl string
	glLink            string
}

var (
	discordHttpClient = http.Client{Timeout: 15 * time.Second}
)

func NewSt(slackWebhookUrl, glLink string) *St {
	return &St{
		discordWebhookUrl: slackWebhookUrl,
		glLink:            glLink,
	}
}

func (c *St) HandleMessage(msgBytes []byte) error {
	msg := MsgSt{}
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return err
	}

	for _, bl := range msg.Backlog {
		rows := make([]string, 0)
		blMsg := map[string]interface{}{}
		if err = json.Unmarshal([]byte(bl.Message), &blMsg); err == nil {
			for k, v := range blMsg {
				rows = append(rows, fmt.Sprintf("       %s: *%v*", k, v))
			}
		} else {
			rows = append(rows, fmt.Sprintf("       message: *%s*", bl.Message))
		}

		if c.glLink != "" {
			rows = append(rows, "<"+c.glLink+"|GrayLog>")
		}
		err = c.discordSend(DiscordMsgSt{
			Username: bl.Fields.ContainerName,
			Content:  strings.Join(rows, "\n"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *St) discordSend(msg DiscordMsgSt) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.discordWebhookUrl, bytes.NewBuffer(msgBytes))
	if err != nil {
		return err
	}

	resp, err := discordHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("bad status code from discord")
	}

	return nil
}
