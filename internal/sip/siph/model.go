package siph

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// common struct
type (
	R struct {
		Result string `xml:"Result"`
	}

	Mata struct {
		CmdType  string `xml:"CmdType"`
		DeviceID string `xml:"DeviceID"`
		SN       string `xml:"SN"`
	}
)

// notify 消息
type (
	// Keepalive xml解析心跳包结构
	Keepalive struct {
		CmdType  string `xml:"CmdType"`
		SN       int    `xml:"SN"`
		DeviceID string `xml:"DeviceID"`
		Status   string `xml:"Status"`
		Info     string `xml:"Info"`
	}

	// AlarmNotify 报警通知
	AlarmNotify struct {
		Mata
		// 报警级别，1为一级警情、2为二级警情、3为三级警情、4为四级警情
		AlarmPriority string `xml:"AlarmPriority"`
		// !-报警方式（必选），取值1为电话报警，2为设备报警，3为短信报警，4为GPS报警
		//5为视频报警，6为设备故障报警，7其他报警-〉
		AlarmMethod string `xml:"AlarmMethod"`
		// 报警时间
		AlarmTime string `xml:"AlarmTime"`
		// 报警内容描述
		AlarmDescription string `xml:"AlarmDescription"`
		// 经纬度信息
		Longitude string `xml:"Longitude"`
		Latitude  string `xml:"Latitude"`
		// 扩展信息
		Info AlarmNotifyInfo `xml:"Info"`
	}

	AlarmNotifyInfo struct {
		// 根据报警方式的不同，该字段的取值也有不同的含义
		// 报警方式为2时，不携带AlarmType为默认的报警设备报警，但如果携带AlarmType取值，那么报警类型如下：1-视频丢失报警；2-设备防拆报警；3-存储设备磁盘满报警；4-设备高温报警；5-设备低温报警。
		//
		// 报警方式为5时，取值如下：1-人工视频报警；2-运动目标检测报警；3-遗留物检测报警；4-物体移除检测报警；5-绊线检测报警6-入侵检测报警；7-逆行检测报警；8-徘徊检测报警；9-流量统计报警；10-密度检测报警；11-视频异常检测报警；12-快速移动报警。
		//
		// 报警方式为6时，取值如下：1-存储设备磁盘故障报警；2-存储设备风扇故障报警。
		AlarmType      string         `xml:"AlarmType"`
		AlarmTypeParam AlarmTypeParam `xml:"AlarmTypeParam"`
	}

	AlarmTypeParam struct {
		//(！一报警类型扩展参数。在人侵检测报警时可携带EventType〉事件类型(/Even- tType〉，事件类型取值：1-进入区域；2-离开区域。-〉
		EventType string `xml:"EventType"`
	}
)

// message消息
type (
	// DeviceInfo 设备信息
	DeviceInfo struct {
		CmdType      string `xml:"CmdType"`
		SN           string `xml:"SN"`
		DeviceID     string `xml:"DeviceID"`
		Result       string `xml:"Result"`
		DeviceName   string `xml:"DeviceName"`
		Manufacturer string `xml:"Manufacturer"`
		Model        string `xml:"Model"`
		Firmware     string `xml:"Firmware"`
	}

	// DeviceCatalogResponse 设备目录响应
	DeviceCatalogResponse struct {
		Name xml.Name `xml:"Response"`
		Mata
		SunNum        string     `xml:"SunNum"`
		DeviceListNum string     `xml:"Num,attr"`
		DeviceList    DeviceList `xml:"DeviceList"`
	}

	// DeviceList 目录列表
	DeviceList struct {
		Name  xml.Name      `xml:"DeviceList"`
		Items []CatalogItem `xml:"Item"`
	}

	// CatalogItem 设备目录中的每一项结构体
	CatalogItem struct {
		XmlName      xml.Name `xml:"Item"`
		DeviceId     string   `xml:"DeviceID"`
		Name         string   `xml:"Name"`
		Manufacturer string   `xml:"Manufacturer"`
		Model        string   `xml:"Model"`
		Owner        string   `xml:"Owner"`
		CivilCode    string   `xml:"CivilCode"`
		Address      string   `xml:"Address"`
		Parental     string   `xml:"Parental"`
		ParentID     string   `xml:"ParentID"`
		SafetyWay    string   `xml:"SafetyWay"`
		RegisterWay  string   `xml:"RegisterWay"`
		Secrecy      string   `xml:"Secrecy"`
		Status       string   `xml:"Status"`
	}

	// DeviceBasicConfigResp 设备基本配置查询返回结构体
	DeviceBasicConfigResp struct {
		Mata
		R
		BasicParam BasicParam `xml:"BasicParam"`
	}

	// BasicParam 设备基本配置Basic配置项结构体
	BasicParam struct {
		Name              string `xml:"Name"`
		DeviceId          string `xml:"DeviceID"`
		SipServerId       string `xml:"SIPServerID"`
		SipServerIp       string `xml:"SIPServerIP"`
		SipServerPort     string `xml:"SIPServerPort"`
		DomainName        string `xml:"DomainName"`
		Expiration        string `xml:"Expiration"`
		Password          string `xml:"Password"`
		HeartBeatInterval int    `xml:"HeartBeatInterval"`
		HeartBeatCount    int    `xml:"HeartBeatCount"`
	}

	DeviceStatus struct {
		Mata
		R
		Reason      string `xml:"Reason"`
		Online      string `xml:"Online"`
		Status      string `xml:"Status"`
		DeviceTime  string `xml:"DeviceTime"`
		Encode      string `xml:"Encode"`
		Record      string `xml:"Record"`
		AlarmStatus Alarm  `xml:"Alarmstatus"`
	}

	Alarm struct {
		Num  int           `xml:"Num,attr"`
		Item []AlarmStatus `xml:"Item"`
	}
	AlarmStatus struct {
		DeviceID   string `xml:"DeviceID"`
		DutyStatus string `xml:"DutyStatus"`
	}
)

func (a AlarmNotify) String() string {
	buf := strings.Builder{}

	// Map报警方式的数字到对应的文本
	alarmMethods := map[string]string{
		"1": "电话报警",
		"2": "设备报警",
		"3": "短信报警",
		"4": "GPS报警",
		"5": "视频报警",
		"6": "设备故障报警",
		"7": "其他报警",
	}

	// Map报警类型的数字到对应的文本
	alarmTypes := map[string]string{
		"1":  "视频丢失报警",
		"2":  "设备防拆报警",
		"3":  "存储设备磁盘满报警",
		"4":  "设备高温报警",
		"5":  "设备低温报警",
		"6":  "人工视频报警",
		"7":  "运动目标检测报警",
		"8":  "遗留物检测报警",
		"9":  "物体移除检测报警",
		"10": "绊线检测报警",
		"11": "入侵检测报警",
		"12": "逆行检测报警",
		"13": "徘徊检测报警",
		"14": "流量统计报警",
		"15": "密度检测报警",
		"16": "视频异常检测报警",
		"17": "快速移动报警",
		"18": "存储设备磁盘故障报警",
		"19": "存储设备风扇故障报警",
	}

	buf.WriteString(fmt.Sprintf("报警级别：%s级警情\n", a.AlarmPriority))

	alarmMethodText, found := alarmMethods[a.AlarmMethod]
	if found {
		buf.WriteString(fmt.Sprintf("报警方式：%s\n", alarmMethodText))
	}

	buf.WriteString(fmt.Sprintf("报警时间：%s\n", a.AlarmTime))
	buf.WriteString(fmt.Sprintf("报警内容描述：%s\n", a.AlarmDescription))
	buf.WriteString(fmt.Sprintf("经纬度信息：%s,%s\n", a.Longitude, a.Latitude))
	buf.WriteString("扩展信息：\n")

	alarmTypeText, found := alarmTypes[a.Info.AlarmType]
	if found {
		buf.WriteString(fmt.Sprintf("\t报警类型：%s\n", alarmTypeText))

		if a.Info.AlarmType == "11" {
			eventTypes := map[string]string{
				"1": "进入区域",
				"2": "离开区域",
			}
			eventTypeText, found := eventTypes[a.Info.AlarmTypeParam.EventType]
			if found {
				buf.WriteString(fmt.Sprintf("\t\t事件类型：%s\n", eventTypeText))
			}
		}
	}

	return buf.String()
}
