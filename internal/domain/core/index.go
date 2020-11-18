package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (c *St) HandleMessage(msgBytes []byte) {
	msg := MsgSt{}
	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		log.Println("Fail to parse json", err)
		return
	}

	for _, bl := range msg.Backlog {
		rows := make([]string, 0)
		blMsg := map[string]interface{}{}
		if err = json.Unmarshal([]byte(bl.Message), &blMsg); err == nil {
			for k, v := range blMsg {
				switch vv := v.(type) {
				case string:
					if len(vv) > 1400 {
						vv = vv[:1400] + "..."
					}
					rows = append(rows, fmt.Sprintf("```%s: %s```", k, vv))
				default:
					rows = append(rows, fmt.Sprintf("```%s: %v```", k, vv))
				}
			}
		} else {
			rows = append(rows, fmt.Sprintf("```message**: %s```", bl.Message))
		}

		if c.glLink != "" {
			rows = append(rows, "<"+c.glLink+">")
		}
		c.discordSend(DiscordMsgSt{
			Username: bl.Fields.ContainerName,
			Content:  strings.Join(rows, ""),
		})
	}
}

func (c *St) discordSend(msg DiscordMsgSt) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Fail to marshal json", err)
		return
	}

	req, err := http.NewRequest("POST", c.discordWebhookUrl, bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Println("Fail to create request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := discordHttpClient.Do(req)
	if err != nil {
		log.Println("Fail to send request", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var bodyText string

		respBodyRaw, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			bodyText = string(respBodyRaw)
		}

		log.Println("bad status code from discord", resp.StatusCode, "body:", bodyText)
	}
}
