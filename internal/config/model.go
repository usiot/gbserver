package config

var Conf struct {
	Redis string `json:"redis"`
	Sip   `json:"sip"`
	Log   `json:"log"`
	Mysql `json:"mysql"`
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

	Log struct {
		Filename   string `json:"filename"`
		MaxSize    int64  `json:"maxSize"`
		MaxAge     int    `json:"maxAge"`
		Level      int    `json:"level"`
		MaxBackups int    `json:"maxBackups"`
	}

	Mysql struct {
		Host     string `json:"host"`
		Port     uint16 `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	}
)
