package proxy

type IProxy interface {
	GetRawUri() string
	SetRtt(rtt int)
}

type Proxy struct {
	RawUri string `koanf,json:"raw_uri"`
	RTT    int    `koanf,json:"rtt"`
}

func (that *Proxy) GetRawUri() string {
	return that.RawUri
}

func (that *Proxy) SetRtt(rtt int) {
	that.RTT = rtt
}
