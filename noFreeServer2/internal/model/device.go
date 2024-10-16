package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	UserID       uint
	Fingerprint  string    //设备指纹
	ExtID        string    //插件ID
	LastActiveAt time.Time //最后活跃时间
	UserAgent    string    //用户代理
	IPAddress    string    //IP地址
}
