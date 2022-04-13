package networking

type Checkpoint struct {
	Block string `json:"block"`
	State string `json:"state"`
	Epoch string `json:"epoch"`
}

type SubscribeOpts struct {
	Endpoints  []string
	StreamURL  string
	Subscriber Subscriber
}

type ValidatorBalanceList struct {
	Data []ValidatorBalance `json:"data"`
}
type ValidatorBalance struct {
	Index   string `json:"index"`
	Balance string `json:"balance"`
}
