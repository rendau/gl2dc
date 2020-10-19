package core

type MsgSt struct {
	Backlog []MsgBacklogSt `json:"backlog"`
	Event   MsgEventSt     `json:"event"`
}

type MsgEventSt struct {
	Fields map[string]interface{} `json:"fields"`
}

type MsgBacklogSt struct {
	Message string             `json:"message"`
	Fields  MsgBacklogFieldsSt `json:"fields"`
}

type MsgBacklogFieldsSt struct {
	ContainerName string `json:"container_name"`
	GlMessageId   string `json:"gl2_message_id"`
}

type DiscordMsgSt struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
