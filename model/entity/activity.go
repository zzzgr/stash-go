package entity

import (
	"gorm.io/gorm"
	"time"
)

type Activity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Code         string  `gorm:"type:varchar(255);not null;comment:'活动code'" json:"code"`
	Name         string  `gorm:"type:varchar(255);not null;comment:'活动名称'" json:"name"`
	Cron         string  `gorm:"type:varchar(40);not null;comment:'cron'" json:"cron"`
	Advance      float64 `gorm:"type:int(11);not null;default:0;comment:'提前时间'" json:"advance"`
	Interval     float64 `gorm:"type:int(11);not null;default:0;comment:'间隔时间'" json:"interval"`
	UrlPattern   string  `gorm:"type:varchar(1000);not null;default:'';comment:'url匹配规则'" json:"urlPattern"`
	QueryPattern string  `gorm:"type:varchar(1000);not null;default:'';comment:'query参数匹配规则,多个用&分隔即可'" json:"queryPattern"`
	Type         int     `gorm:"type:tinyint(1);not null;default:1;comment:'分号的类型 1cookie 2header'" json:"type"`
	Field        string  `gorm:"type:varchar(64);not null;comment:default:'';'分号的依据 多个用英文逗号分割'" json:"field"`

	AlertAhead int `gorm:"type:int(11);not null;default:0;comment:'提前通知时间 分钟'" json:"alertAhead"`

	Status int `gorm:"type:tinyint(1);not null;default:1;comment:'状态 1启用 2禁用'" json:"status"`
}
