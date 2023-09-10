package dao

import (
	"context"
	"time"

	"github.com/usiot/gbserver/internal/dao/internal"
)

//go:generate sqlstruct -types=DbDevice,DbPtrDevice -file=device.go
type (
	DbDevice struct {
		DeviceId          string `json:"deviceId"          db:"DEVICE_ID"`          // 设备序列号
		Domain            string `json:"domain"            db:"DOMAIN"`             // 设备的 sip 域名
		Name              string `json:"name"              db:"NAME"`               // 设备名
		Manufacturer      string `json:"manufacturer"      db:"MANUFACTURER"`       // 制造厂商
		Model             string `json:"model"             db:"MODEL"`              // 设备型号
		Firmware          string `json:"firmware"          db:"FIRMWARE"`           // 固件版本
		Transport         string `json:"transport"         db:"TRANSPORT"`          // 传输模式
		Status            uint8  `json:"status"            db:"STATUS"`             // 在离线状态，0：离线，1：在线
		Ip                string `json:"ip"                db:"IP"`                 // ip 地址
		Port              uint16 `json:"port"              db:"PORT"`               // 传输端口
		Expires           int64  `json:"expires"           db:"EXPIRES"`            // 心跳过期时间
		RegisterTime      int64  `json:"registerTime"      db:"REGISTER_TIME"`      // 注册时间
		Keepalive         int64  `json:"keepalive"         db:"KEEPALIVE"`          // 上次心跳时间
		HeartbeatCount    uint   `json:"heartbeatCount"    db:"HEARTBEAT_COUNT"`    // 心跳超时次数，范围值: 3-255
		HeartbeatInterval uint   `json:"heartbeatInterval" db:"HEARTBEAT_INTERVAL"` // 心跳间隔时间，范围值: 5-255

		CreatedAt int64 `json:"createdAt" db:"CREATED_AT"` // 创建时间
		UpdatedAt int64 `json:"updatedAt" db:"UPDATED_AT"` // 更新时间
	}

	DbPtrDevice struct {
		DeviceId          string  `json:"deviceId"          db:"DEVICE_ID"`          // 设备序列号
		Domain            *string `json:"domain"            db:"DOMAIN"`             // 设备的 sip 域名
		Name              *string `json:"name"              db:"NAME"`               // 设备名
		Manufacturer      *string `json:"manufacturer"      db:"MANUFACTURER"`       // 制造厂商
		Model             *string `json:"model"             db:"MODEL"`              // 设备型号
		Firmware          *string `json:"firmware"          db:"FIRMWARE"`           // 固件版本
		Transport         *string `json:"transport"         db:"TRANSPORT"`          // 传输模式
		Status            *uint8  `json:"status"            db:"STATUS"`             // 设备在离线状态，0：离线，1：在线
		Ip                *string `json:"ip"                db:"IP"`                 // ip 地址
		Port              *uint16 `json:"port"              db:"PORT"`               // 传输端口
		Expires           *int64  `json:"expires"           db:"EXPIRES"`            // 心跳过期时间
		RegisterTime      *int64  `json:"registerTime"      db:"REGISTER_TIME"`      // 注册时间
		Keepalive         *int64  `json:"keepalive"         db:"KEEPALIVE"`          // 上次心跳时间
		HeartbeatCount    *uint   `json:"heartbeatCount"    db:"HEARTBEAT_COUNT"`    // 心跳超时次数，范围值: 3-255
		HeartbeatInterval *uint   `json:"heartbeatInterval" db:"HEARTBEAT_INTERVAL"` // 心跳间隔时间，范围值: 5-255

		CreatedAt *int64 `json:"createdAt" db:"CREATED_AT"` // 创建时间
		UpdatedAt *int64 `json:"updatedAt" db:"UPDATED_AT"` // 更新时间
	}
)

var (
	Offline uint8 = 0
	Online  uint8 = 1
)

var (
	_ SQLFns    = &DbDevice{}
	_ SQLPtrFns = &DbPtrDevice{}
)

func GetDeviceById(ctx context.Context, did string) (*DbDevice, error) {
	var device DbDevice
	err := internal.BeginSelect(&device).
		From(TableDevice).
		Where(map[string]interface{}{"DEVICE_ID": did}).
		GetOne(ctx, db, &device)
	return &device, err
}

func ChangeDeviceStatus(ctx context.Context, did string, status uint8) error {
	return Update(ctx, TableDevice, &DbPtrDevice{
		DeviceId: did,
		Status:   &status,
	})
}

func UpdateDevice(ctx context.Context, dev *DbDevice) error {
	return Update(ctx, TableDevice, &DbPtrDevice{
		DeviceId:          dev.DeviceId,
		Domain:            &dev.Domain,
		Name:              &dev.Name,
		Manufacturer:      &dev.Manufacturer,
		Model:             &dev.Model,
		Firmware:          &dev.Firmware,
		Transport:         &dev.Transport,
		Status:            &dev.Status,
		Ip:                &dev.Ip,
		Port:              &dev.Port,
		Expires:           &dev.Expires,
		RegisterTime:      &dev.RegisterTime,
		Keepalive:         &dev.Keepalive,
		HeartbeatCount:    &dev.HeartbeatCount,
		HeartbeatInterval: &dev.HeartbeatInterval,
		CreatedAt:         &dev.CreatedAt,
		UpdatedAt:         &dev.UpdatedAt,
	})
}

func (i *DbDevice) SQLValues() internal.SQLPairs {
	return []SQLPair{
		{K: "`DEVICE_ID`", V: i.DeviceId},
		{K: "`DOMAIN`", V: i.Domain},
		{K: "`NAME`", V: i.Name},
		{K: "`MANUFACTURER`", V: i.Manufacturer},
		{K: "`MODEL`", V: i.Model},
		{K: "`FIRMWARE`", V: i.Firmware},
		{K: "`TRANSPORT`", V: i.Transport},
		{K: "`OFFLINE`", V: i.Status},
		{K: "`IP`", V: i.Ip},
		{K: "`PORT`", V: i.Port},
		{K: "`EXPIRES`", V: i.Expires},
		{K: "`REGISTER_TIME`", V: i.RegisterTime},
		{K: "`KEEPALIVE`", V: i.Keepalive},
		{K: "`HEARTBEAT_COUNT`", V: i.HeartbeatCount},
		{K: "`HEARTBEAT_INTERVAL`", V: i.HeartbeatInterval},
		{K: "`CREATED_AT`", V: i.CreatedAt},
		{K: "`UPDATED_AT`", V: i.UpdatedAt},
	}
}

func (i *DbPtrDevice) SQLValues() SQLPairs {
	return []SQLPair{
		{K: "`DEVICE_ID`", V: i.DeviceId},
		{K: "`DOMAIN`", V: i.Domain},
		{K: "`NAME`", V: i.Name},
		{K: "`MANUFACTURER`", V: i.Manufacturer},
		{K: "`MODEL`", V: i.Model},
		{K: "`FIRMWARE`", V: i.Firmware},
		{K: "`TRANSPORT`", V: i.Transport},
		{K: "`OFFLINE`", V: i.Status},
		{K: "`IP`", V: i.Ip},
		{K: "`PORT`", V: i.Port},
		{K: "`EXPIRES`", V: i.Expires},
		{K: "`REGISTER_TIME`", V: i.RegisterTime},
		{K: "`KEEPALIVE`", V: i.Keepalive},
		{K: "`HEARTBEAT_COUNT`", V: i.HeartbeatCount},
		{K: "`HEARTBEAT_INTERVAL`", V: i.HeartbeatInterval},
		{K: "`CREATED_AT`", V: i.CreatedAt},
		{K: "`UPDATED_AT`", V: i.UpdatedAt},
	}
}

func (i *DbPtrDevice) SQLPtrNotPtrValues() SQLPairs {
	return []SQLPair{
		{K: "`DEVICE_ID`", V: i.DeviceId},
	}
}

func (i *DbPtrDevice) SQLPtrNotNilValues() SQLPairs {
	vals := []SQLPair{}
	if i.Domain != nil {
		vals = append(vals, SQLPair{K: "`DOMAIN`", V: i.Domain})
	}
	if i.Name != nil {
		vals = append(vals, SQLPair{K: "`NAME`", V: i.Name})
	}
	if i.Manufacturer != nil {
		vals = append(vals, SQLPair{K: "`MANUFACTURER`", V: i.Manufacturer})
	}
	if i.Model != nil {
		vals = append(vals, SQLPair{K: "`MODEL`", V: i.Model})
	}
	if i.Firmware != nil {
		vals = append(vals, SQLPair{K: "`FIRMWARE`", V: i.Firmware})
	}
	if i.Transport != nil {
		vals = append(vals, SQLPair{K: "`TRANSPORT`", V: i.Transport})
	}
	if i.Status != nil {
		vals = append(vals, SQLPair{K: "`OFFLINE`", V: i.Status})
	}
	if i.Ip != nil {
		vals = append(vals, SQLPair{K: "`IP`", V: i.Ip})
	}
	if i.Port != nil {
		vals = append(vals, SQLPair{K: "`PORT`", V: i.Port})
	}
	if i.Expires != nil {
		vals = append(vals, SQLPair{K: "`EXPIRES`", V: i.Expires})
	}
	if i.RegisterTime != nil {
		vals = append(vals, SQLPair{K: "`REGISTER_TIME`", V: i.RegisterTime})
	}
	if i.Keepalive != nil {
		vals = append(vals, SQLPair{K: "`KEEPALIVE`", V: i.Keepalive})
	}
	if i.HeartbeatCount != nil {
		vals = append(vals, SQLPair{K: "`HEARTBEAT_COUNT`", V: i.HeartbeatCount})
	}
	if i.HeartbeatInterval != nil {
		vals = append(vals, SQLPair{K: "`HEARTBEAT_INTERVAL`", V: i.HeartbeatInterval})
	}
	if i.CreatedAt != nil {
		vals = append(vals, SQLPair{K: "`CREATED_AT`", V: i.CreatedAt})
	}
	if i.UpdatedAt != nil {
		vals = append(vals, SQLPair{K: "`UPDATED_AT`", V: i.UpdatedAt})
	}
	return vals
}

func (i *DbPtrDevice) SQLFixPtr() {

	if i.Domain != nil {
		i.Domain = new(string)
	}
	if i.Name != nil {
		i.Name = new(string)
	}
	if i.Manufacturer != nil {
		i.Manufacturer = new(string)
	}
	if i.Model != nil {
		i.Model = new(string)
	}
	if i.Firmware != nil {
		i.Firmware = new(string)
	}
	if i.Transport != nil {
		i.Transport = new(string)
	}
	if i.Status != nil {
		i.Status = new(uint8)
	}
	if i.Ip != nil {
		i.Ip = new(string)
	}
	if i.Port != nil {
		i.Port = new(uint16)
	}
	if i.Expires != nil {
		i.Expires = new(int64)
	}
	if i.RegisterTime != nil {
		i.RegisterTime = new(int64)
	}
	if i.Keepalive != nil {
		i.Keepalive = new(int64)
	}
	if i.HeartbeatCount != nil {
		i.HeartbeatCount = new(uint)
	}
	if i.HeartbeatInterval != nil {
		i.HeartbeatInterval = new(uint)
	}
	if i.CreatedAt != nil {
		i.CreatedAt = new(int64)
	}
	if i.UpdatedAt != nil {
		i.UpdatedAt = new(int64)
	}
}

func (i *DbDevice) SetUpdateTime(now time.Time) {
	unixMill := now.UnixMilli()
	i.UpdatedAt = unixMill
}

func (i *DbDevice) SetCreateTime(now time.Time) {
	unixMill := now.UnixMilli()
	i.CreatedAt = unixMill
	i.UpdatedAt = unixMill
}

func (i *DbPtrDevice) SetUpdateTime(now time.Time) {
	unixMill := now.UnixMilli()
	i.UpdatedAt = &unixMill
}

func (i *DbPtrDevice) SetCreateTime(now time.Time) {
	unixMill := now.UnixMilli()
	i.CreatedAt = &unixMill
	i.UpdatedAt = &unixMill
}
