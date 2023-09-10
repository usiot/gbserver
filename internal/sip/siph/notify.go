package siph

import (
	"context"
	"net/http"
	"time"

	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

func KeepaliveNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	keepalive := Keepalive{}
	err := UnmarshalReq(req.Body(), &keepalive)
	if err != nil {
		logger.Debug("keepalive 消息解析xml失败：%s", err)
		return
	}
	device, ok := DeviceFromRequest(req)
	if !ok {
		return
	}

	ctx := context.WithValue(context.Background(), util.CtxDid, device.DeviceId)

	// TODO:
	// 获取设备信息
	device, err = dao.GetDeviceById(ctx, device.DeviceId)
	if err != nil {
		resp := sip.NewResponseFromRequest("", req, http.StatusNotFound, "device "+device.DeviceId+"not found", "")
		logger.Debug("{%s}设备不存在\n%s", device.DeviceId, resp)
		_ = tx.Respond(resp)
		return
	}

	// 更新心跳时间
	unixMill := time.Now().UnixMilli()
	dao.Update(ctx, dao.TableDevice, &dao.DbPtrDevice{
		DeviceId:  device.DeviceId,
		Status:    &dao.Online,
		Keepalive: &unixMill,
	})
	// TODO: 维护定时任务

	resp := sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), "")
	_ = tx.Respond(resp)
}

func AlarmNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	alarm := AlarmNotify{}
	err := UnmarshalReq(req.Body(), &alarm)
	if err != nil {
		logger.Debug("alarm 消息解析xml失败：%s", err)
		return
	}

	logger.Debug("alarm notify : \n%s", alarm.String())

	_ = responseAck(tx, req)
}

func MobilePositionNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	// 自行扩展

	// logger.Debug("mobile position notify: \n%s", req.Body())

	_ = responseAck(tx, req)
}
