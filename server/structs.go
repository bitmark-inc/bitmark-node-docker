package server

type Counters struct {
	Pending  int `json:"pending"`
	Verified int `json:"verified"`
}
type BlockInfo struct {
	LRCount `json:"count"`
	Hash    string `json:"hash"`
}
type LRCount struct {
	Local  uint64 `json:"local"`
	Remote uint64 `json:"remote"`
}

type PeerCounts struct {
	Incoming uint64 `json:"incoming"`
	Outgoing uint64 `json:"outgoing"`
}

type DetailReply struct {
	Chain               string     `json:"chain"`
	Mode                string     `json:"mode"`
	Block               BlockInfo  `json:"block"`
	RPCs                uint64     `json:"rpcs"`
	Peers               PeerCounts `json:"peers"`
	TransactionCounters Counters   `json:"transactionCounters"`
	Difficulty          float64    `json:"difficulty"`
	Hashrate            float64    `json:"hashrate,omitempty"`
	Version             string     `json:"version"`
	Uptime              string     `json:"uptime"`
	PublicKey           string     `json:"publicKey"`
}
