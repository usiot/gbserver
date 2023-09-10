package sip

import (
	"context"
	"fmt"
	"time"

	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"github.com/usiot/gbserver/internal/cache"
	"github.com/usiot/gbserver/internal/config"
	"github.com/usiot/gbserver/internal/dao"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

type Server struct {
	host    string
	network string
	sip     gosip.Server
	cfg     *config.Sip
}

const (
	letterBytes    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	contentTypeXML = "Application/MANSCDP+xml"
	contentTypeSDP = "APPLICATION/SDP"
)

var sipSrv Server

func Init(cfg *config.Sip) {
	sipSrv = Server{
		host:    fmt.Sprintf("%s:%d", cfg.Host, cfg.SipPort),
		network: cfg.Network,
		sip: gosip.NewServer(
			gosip.ServerConfig{UserAgent: cfg.UserAgent},
			nil,
			nil,
			log.NewDefaultLogrusLogger(),
		),
		cfg: cfg,
	}
	sipSrv.Register()

	util.Go(func() {
		err := sipSrv.ListenTCP()
		if err != nil {
			logger.Fatal(err.Error())
		}
	})
	util.Go(func() {
		err := sipSrv.ListenUDP()
		if err != nil {
			logger.Fatal(err.Error())
		}
	})
}

func (s *Server) Register() {
	s.register(sip.REGISTER, RegisterHandler)
	s.register(sip.MESSAGE, MessageHandler)
}

func (s *Server) register(method sip.RequestMethod, handler gosip.RequestHandler) error {
	return s.sip.OnRequest(method, handler)
}

func (s *Server) ListenTCP() error {
	logger.Info("gb server listen tcp: %s", s.host)
	return s.sip.Listen("tcp", s.host, nil)
}

func (s *Server) ListenUDP() error {
	logger.Info("gb server listen udp: %s", s.host)
	return s.sip.Listen("udp", s.host, nil)
}

func (s *Server) Shutdown() error {
	s.sip.Shutdown()
	logger.Info("gb server shutdown...")
	return nil
}

func (s *Server) QueryDeviceInfo(ctx context.Context, dev *dao.DbDevice) {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Query>
  <CmdType>DeviceInfo</CmdType>
  <SN>701385</SN>
  <DeviceID>%s</DeviceID>
</Query>`, dev.DeviceId)

	req := s.createMessageRequest(ctx, dev, body)
	logger.Debug("查询设备信息请求：\n%s", req)
	trans, _ := s.sip.Request(req)
	if trans != nil {
		resp := <-trans.Responses()
		logger.Debug("收到设备查询响应：\n%s", resp)
	}
	s.QueryDeviceCatalog(ctx, dev)
}

func (s *Server) QueryDeviceCatalog(ctx context.Context, dev *dao.DbDevice) {
	body := `<?xml version="1.0" encoding="utf-8"?>
<Query>
  <CmdType>Catalog</CmdType>
  <SN>98760</SN>
  <DeviceID>44010200491118000001</DeviceID>
</Query>`
	req := s.createMessageRequest(ctx, dev, body)
	logger.Debug("查询设备目录查询请求：\n%s", req)
	trans, _ := s.sip.Request(req)
	if trans != nil {
		resp := <-trans.Responses()
		logger.Debug("收到设备目录查询响应：\n%s", resp)
	}
}

func (s *Server) QueryDeviceStatus(ctx context.Context, dev *dao.DbDevice) {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Query>
  <CmdType>DeviceStatus</CmdType>
  <SN>0</SN>
  <DeviceID>%s</DeviceID>
</Query>`, dev.DeviceId)
	req := s.createMessageRequest(ctx, dev, body)
	logger.Debug("查询设备状态查询请求：\n%s", req)
	trans, _ := s.sip.Request(req)
	if trans != nil {
		resp := <-trans.Responses()
		logger.Debug("收到设备状态查询响应：\n%s", resp)
	}
}

func (s *Server) createMessageRequest(ctx context.Context, d *dao.DbDevice, body string) sip.Request {
	rb := sip.NewRequestBuilder()

	params := sip.NewParams()
	params.Add("tag", sip.String{Str: string(util.NanoId(32))})

	rb.SetFrom(&sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: s.cfg.SipId},
			FHost: s.cfg.SipDomain,
		},
		Params: params,
	})

	port := sip.Port(d.Port)
	to := &sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: d.DeviceId},
			FHost: d.Ip,
			FPort: &port,
		},
	}
	rb.SetTo(to)
	rb.SetRecipient(to.Uri)
	rb.AddVia(s.newVia(d.Transport))
	contentType := sip.ContentType(contentTypeXML)
	rb.SetContentType(&contentType)
	rb.SetMethod(sip.MESSAGE)
	userAgent := sip.UserAgentHeader("gbserver-usiot")
	rb.SetUserAgent(&userAgent)
	rb.SetBody(body)

	ceq, err := cache.GetCeq(ctx)
	if err == nil {
		rb.SetSeqNo(uint(ceq))
	}
	req, _ := rb.Build()
	return req
}

func (s Server) newVia(transport string) *sip.ViaHop {
	p := sip.Port(s.cfg.SipPort)

	params := sip.NewParams()
	params.Add("branch", sip.String{Str: fmt.Sprintf("%s%d", "z9hG4bK", time.Now().UnixMilli())})

	return &sip.ViaHop{
		ProtocolName:    "SIP",
		ProtocolVersion: "2.0",
		Transport:       transport,
		Host:            s.cfg.SipAddress,
		Port:            &p,
		Params:          params,
	}
}
