package sip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/sip/siph"
	"github.com/usiot/gbserver/internal/util"
)

func RegisterHandler(req sip.Request, tx sip.ServerTransaction) {
	logger.Debug("收到register请求\n%s", printRequest(req))
	// 判断是否存在 Authorization 字段
	if headers := req.GetHeaders("Authorization"); len(headers) > 0 { // 存在 Authorization 字段
		fromRequest, ok := siph.DeviceFromRequest(req)
		if !ok {
			return
		}
		ctx := context.WithValue(context.Background(), util.CtxDid, fromRequest.DeviceId)

		offlineFlag := false
		device, err := dao.GetDeviceById(ctx, fromRequest.DeviceId)
		if err != nil {
			logger.Debug("not found from device from database")
			device = fromRequest
		}
		h := req.GetHeaders(ExpiresHeader)
		if len(h) == 0 {
			logger.Error("not found expires header from request", req)
			return
		}
		expires, ok := h[0].(*sip.Expires)
		if !ok {
			logger.Error("not found expires header from request", req)
			return
		}
		// 如果v=0，则代表该请求是注销请求
		if expires.Equals(new(sip.Expires)) {
			logger.Debug("expires值为0,该请求是注销请求")
			offlineFlag = true
		}
		if expires != nil {
			device.Expires = int64(*expires)
		}
		logger.Info("设备信息:  %+v\n", device)

		// 发送OK信息
		resp := sip.NewResponseFromRequest("", req, http.StatusOK, "ok", "")
		logger.Debug("发送OK信息\n%s", resp)
		tx.Respond(resp)

		if offlineFlag { // 注销请求
			device.Status = dao.Offline
			dao.UpdateDevice(ctx, device)
			// TODO: 注销定时任务
		} else {
			// 注册请求
			device.Status = dao.Online
			dao.UpdateDevice(ctx, device)
			go sipSrv.QueryDeviceInfo(ctx, device)
		}
		return
	}

	// 没有存在 Authorization 头部字段
	resp := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "")
	// 添加 WWW-Authenticate 头
	wwwHeader := &sip.GenericHeader{
		HeaderName: WwwHeader,
		Contents: fmt.Sprintf(
			`Digest nonce="%s", algorithm=%s, realm="%s", qop="auth"`,
			"44010200491118000001",
			DefaultAlgorithm,
			util.NanoId(32),
		),
	}
	resp.AppendHeader(wwwHeader)
	logger.Debug("没有Authorization头部信息，生成WWW-Authenticate头部返回：\n%s", resp)
	_ = tx.Respond(resp)
}
