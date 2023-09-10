package siph

import (
	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/logger"
)

func DeviceConfigQueryHandler(req sip.Request, tx sip.ServerTransaction) {
	// logger.Debug("获取到的configDownload消息：\n%s", req.Body())
	defer func() {
		_ = responseAck(tx, req)
	}()

	var cfg DeviceBasicConfigResp
	err := UnmarshalReq(req.Body(), &cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if cfg.R.Result != "OK" {
		return
	}

	// syn.HasSyncTask(fmt.Sprintf("%s_%s", syn.KeyControlDeviceConfigQuery, cfg.DeviceID.DeviceID), func(e *syn.Entity) {
	// 	e.Ok(cfg)
	// })

	// _ = storage.updateDeviceBasicConfig(*cfg)
}
