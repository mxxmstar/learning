package domain

import "time"

// 登录上下文信息
type LoginContext struct {
	DeviceId      string
	IPAddress     string
	UserAgent     string // 客户端设备类型，版本等信息
	ClientVersion string // 客户端版本
	OSInfo        string
	Location      string
}

// 用户固有属性
type User struct {
	Id       uint64
	Username string
	Email    string
	Password string
	CTime    time.Time
}
