package sip

import (
	"net/http"

	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/sip/siph"
)

func MessageHandler(req sip.Request, tx sip.ServerTransaction) {
	logger.Debug("收到MESSAGE请求\n%s", printRequest(req))
	if l, ok := req.ContentLength(); !ok || l.Equals(0) {
		resp := sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), "")
		logger.Debug("该MESSAGE消息的消息体长度为0，返回OK\n%s", resp)
		_ = tx.Respond(resp)
	}

	body := req.Body()
	cmdType, err := siph.GetCmdTypeFromXML(body)
	logger.Debug("解析出的命令：%s", cmdType)
	if err != nil {
		return
	}

	switch cmdType {
	case "Notify:Keepalive":
		siph.KeepaliveNotifyHandler(req, tx)
	case "Notify:Alarm":
		siph.AlarmNotifyHandler(req, tx)
		go sipSrv.QueryDeviceStatus(nil, &dao.DbDevice{
			DeviceId:  "34020000001320000001",
			Port:      5060,
			Ip:        "192.168.124.26",
			Transport: "UDP",
		})
	case "Notify:MobilePosition":
		siph.MobilePositionNotifyHandler(req, tx)
	case "Response:DeviceInfo": // 查询设备信息响应
		siph.DeviceInfoHandler(req, tx)
	case "Response:DeviceConfig": // 设备配置请求应答
		siph.DeviceConfigResponseHandler(req, tx)
	case "Response:Catalog": // 查询设备目录信息响应
		siph.CatalogHandler(req, tx)
	case "Response:DeviceStatus": // 查询设备状态信息响应
		siph.DeviceStatusHandler(req, tx)
	case "Response:ConfigDownload": // 查询设备配置信息响应
		siph.DeviceConfigQueryHandler(req, tx)
	case "Response:Alarm": // 发起报警订阅信息响应
		siph.SubscribeAlarmResponseHandler(req, tx)
	case "Response:MobilePosition": // 发起设备移动位置订阅响应
		siph.SubscribeMobilePositionResponseHandler(req, tx)
	default:
		logger.Warn("不支持的Message方法实现[%s]", cmdType)
	}
}
