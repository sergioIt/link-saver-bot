package telegram

type UpdateResponse struct {
	OK     bool      `json:"ok"`
	Result []Updates `json:"result"`
}

type Updates struct {
	ID      int    `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From `json:"from"`
	Caht Caht `json:"chat"`
}

type From struct {
	UserName string `json:"username"`
}

type Caht struct {
	ID int `json:"id"`
}

