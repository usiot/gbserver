package siph

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func DeviceFromRequest(req sip.Request) (*dao.DbDevice, bool) {
	d := &dao.DbDevice{}

	from, ok := req.From()
	if !ok {
		logger.Debug("从请求中无法解析from头部信息: %s", req.String())
		return d, false
	}

	if from.Address == nil {
		logger.Debug("从请求中无法解析from头address部分信息: %s", req.String())
		return d, false
	}

	if from.Address.User() == nil {
		logger.Debug("从请求中无法解析from头user部分信息: %s", req.String())
		return d, false
	}

	d.DeviceId = from.Address.User().String()
	d.Domain = from.Address.Host()
	via, ok := req.ViaHop()
	if !ok {
		logger.Debug("从请求中无法解析出via头部信息: %s", via.String())
		return d, false
	}

	d.Ip = via.Host
	if via.Port != nil {
		d.Port = uint16(*via.Port)
	}
	d.Transport = via.Transport

	logger.Debug("从请求中解析出的设备信息: %v", d)

	return d, true
}

// GetCmdTypeFromXML 根据body获取XML配置文件中的根元素
func GetCmdTypeFromXML(body string) (key string, err error) {
	decoder := xml.NewDecoder(strings.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "GB2312" {
			return transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder()), nil
		}
		return input, nil
	}

	var (
		isRoot, isCmdType = false, false
		root, cmdType     string
	)

re:
	for t, err := decoder.Token(); err == nil || err == io.EOF; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			if !isRoot {
				root = token.Name.Local
				isRoot = true
			}
			if token.Name.Local == "CmdType" {
				isCmdType = true
			}
		case xml.CharData:
			if isCmdType {
				cmdType = string(token)
				break re
			}
		default:
		}
	}

	key = fmt.Sprintf("%s:%s", root, cmdType)
	return
}

func GetResultFromXML(body string) string {
	_, v := GetSpecificFromXML(body, "Result")
	return v
}

// 在body查询指定key的value，然后返回
func GetSpecificFromXML(body, key string) (k, v string) {
	decoder := xml.NewDecoder(strings.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "GB2312" {
			return transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder()), nil
		}
		return input, nil
	}

	isSpecific := false

re:
	for t, err := decoder.Token(); err == nil || err == io.EOF; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == key {
				isSpecific = true
			}
		case xml.CharData:
			if isSpecific {
				v = string(token)
				break re
			}
		default:
		}
	}
	return key, v
}

func responseAck(transaction sip.ServerTransaction, request sip.Request) error {
	resp := sip.NewResponseFromRequest("", request, sip.StatusCode(http.StatusOK),
		http.StatusText(http.StatusOK), "")
	logger.Debug("发送 ack 响应\n%s", resp)
	err := transaction.Respond(resp)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func UnmarshalReq(body string, v interface{}) error {
	decoder := xml.NewDecoder(strings.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "GB2312" {
			return transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder()), nil
		}
		return input, nil
	}

	return decoder.Decode(v)
}
