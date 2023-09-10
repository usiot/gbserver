package siph

import (
	"context"

	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

type deviceInfo struct {
	CmdType      string `xml:"CmdType"`
	SN           string `xml:"SN"`
	DeviceID     string `xml:"DeviceID"`
	Result       string `xml:"Result"`
	DeviceName   string `xml:"DeviceName"`
	Manufacturer string `xml:"Manufacturer"`
	Model        string `xml:"Model"`
	Firmware     string `xml:"Firmware"`
}

func DeviceInfoHandler(req sip.Request, tx sip.ServerTransaction) {
	d := deviceInfo{}
	if err := UnmarshalReq(req.Body(), &d); err != nil {
		logger.Error("解析deviceInfo响应包出错 %s", err)
		return
	}

	if d.Result != "OK" {
		logger.Error("查询设备信息请求结果为%s，请检查", d.Result)
		return
	}

	ctx := context.WithValue(context.Background(), util.CtxDid, d.DeviceID)

	err := dao.Update(ctx, dao.TableDevice, &dao.DbPtrDevice{
		DeviceId:     d.DeviceID,
		Name:         &d.DeviceName,
		Manufacturer: &d.Manufacturer,
		Model:        &d.Model,
		Firmware:     &d.Firmware,
	})

	if err != nil {
		logger.Error("更新设备信息失败", err)
		return
	}
	responseAck(tx, req)
}

func DeviceConfigResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := GetResultFromXML(req.Body())
	if r == "" {
		logger.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		logger.Error("发送修改配置请求失败，请检查")
	} else {
		logger.Debug("发送修改配置请求成功")
	}
}

func CatalogHandler(req sip.Request, tx sip.ServerTransaction) {
	defer func() {
		_ = responseAck(tx, req)
	}()

	catalog := DeviceCatalogResponse{}

	err := UnmarshalReq(req.Body(), &catalog)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug("catalog: \n%+#v", catalog)
}

func DeviceStatusHandler(req sip.Request, tx sip.ServerTransaction) {
	defer func() {
		_ = responseAck(tx, req)
	}()

	var ds DeviceStatus
	err := UnmarshalReq(req.Body(), &ds)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("device status: \n%+#v", ds)
}

func SubscribeAlarmResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := GetResultFromXML(req.Body())
	if r == "" {
		logger.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		logger.Error("订阅报警信息失败，请检查")
	} else {
		logger.Debug("订阅报警信息成功")
	}

	_ = responseAck(tx, req)
}

func SubscribeMobilePositionResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := GetResultFromXML(req.Body())
	if r == "" {
		logger.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		logger.Error("订阅设备移动位置信息失败，请检查")
	} else {
		logger.Debug("订阅设备移动位置信息信息成功")
	}

	_ = responseAck(tx, req)
}
