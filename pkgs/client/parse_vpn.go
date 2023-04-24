package client

type InBound struct {
	Listen string `json:"listen"`
	Port   int    `json:"port"`
	Raw    string `json:"raw"`
}
