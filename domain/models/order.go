package models

import "gorm.io/gorm"

// 按照proto文件定义的结构体
type Order struct {
	gorm.Model
	OrderId      string `gorm:"column:order_id;unique" json:"order_id"`
	OrderVersion int64  `gorm:"column:order_version" json:"order_version"`
	UserId       int64  `gorm:"column:user_id" json:"user_id"`
	// 用于存储订单数据的json字符串 默认使用varchar 类型，是OrderInfo slice的json字符串
	OrderData string `gorm:"column:order_data" json:"order_data"`
	Status    int8   `gorm:"column:status" json:"status"` // 是否支付
}
