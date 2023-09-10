package config

var Conf struct {
	Sip    `json:"sip"`
	Logger `json:"logger"`
}

type (
	Sip struct {
		SipId      string `json:"sipId"`
		SipDomain  string `json:"sipDomain"`
		Host       string `json:"host"`
		SipPort    uint16 `json:"sipPort"`
		SipAddress string `json:"sipAddress"`
		Network    string `json:"network"`
		UserAgent  string `json:"userAgent"`
	}

	Logger struct {
		Filename   string `json:"filename"`
		MaxSize    int64  `json:"maxSize"`
		MaxAge     int    `json:"maxAge"`
		Level      int    `json:"level"`
		MaxBackups int    `json:"maxBackups"`
	}
)
