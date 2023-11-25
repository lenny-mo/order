package models

import "time"

// OrderItem 订单信息
type OrderItem struct {
	// OrderId 是Order表的外键
	OrderId   int64     `json:"order_id" gorm:"order_id;not_null"`
	ID        int64     `json:"id" gorm:"id;auto_increment;primary_key;not_null"`
	UserId    int64     `json:"user_id" gorm:"user_id;not_null"`
	SKUId     int64     `json:"sku_id" gorm:"sku_id;not_null"`
	Count     int32     `json:"count" gorm:"count;not_null"`
	Timestamp time.Time `json:"timestamp" gorm:"timestamp;not_null"`
	// 外键策略：当Order表的OrderId字段更新时，OrderItem表的OrderId字段也更新
	// 当Order表的OrderId字段删除时，OrderItem表的OrderId字段也删除
	Order Order `json:"order" gorm:"foreignkey:OrderId;references:OrderId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
