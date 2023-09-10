package sip

import (
	"bytes"
	"io"
	"strings"

	"github.com/ghettovoice/gosip/sip"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	DefaultAlgorithm = "MD5"
	WwwHeader        = "WWW-Authenticate"
	ExpiresHeader    = "Expires"
)

func printRequest(req sip.Request) string {
	var buffer bytes.Buffer

	buffer.WriteString(req.StartLine() + "\r\n")
	for _, v := range req.Headers() {
		buffer.WriteString(v.String() + "\r\n")
	}
	buffer.WriteString("\r\n")
	body := req.Body()
	switch {
	case strings.Contains(body, ` encoding="GB2312"`):
		utf8Text, err := GB2312ToUTF8(body)
		if err == nil {
			body = utf8Text
		}
	}
	buffer.WriteString(body)

	return buffer.String()
}

func GB2312ToUTF8(gb2312Text string) (string, error) {
	reader := transform.NewReader(strings.NewReader(gb2312Text), simplifiedchinese.GB18030.NewDecoder())
	utf8Text, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(utf8Text), nil
}
