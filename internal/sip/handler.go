package sip

import (
	"bytes"

	"github.com/ghettovoice/gosip/sip"
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
	buffer.WriteString(body)

	return buffer.String()
}
