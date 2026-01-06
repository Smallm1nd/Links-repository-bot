package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMassage `json:"message"`
}

type IncomingMassage struct {
	Text string `json:"string"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int `json:"id"`
}
type From struct {
	UserName string `json:"username"`
}
