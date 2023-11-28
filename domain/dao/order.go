package dao

import (
	"encoding/json"

	"github.com/lenny-mo/order/domain/models"
	"gorm.io/gorm"
)

type OrderDAOInterface interface {
	CreateOrder(order *models.Order) (int64, error)
	UpdateOrder(order *models.Order) (int64, error)
	GetOrderById(orderId int64) (*models.Order, error)
	CreateOrderItem(orderItem *models.OrderItem) (int64, error)
	// 不是很清楚是否需要更新OrderItem，因为有外键约束，当Order表的OrderId字段更新时，OrderItem表的OrderId字段也更新
}

type OrderDAO struct {
	db *gorm.DB
}

func NewOrderDAO(db *gorm.DB) OrderDAOInterface {
	return &OrderDAO{
		db: db,
	}
}

//	CreateOrder 创建订单
//
// 需要同时插入OrderItem表和Order表
func (o *OrderDAO) CreateOrder(order *models.Order) (rowAffected int64, err error) {
	// 开启事务
	trans := o.db.Begin()

	// 处理不可预见的panic or 事务提交失败
	defer func() {
		if r := recover(); r != nil || trans.Error != nil {
			trans.Rollback()
		} else {
			// 判断提交是否成功
			if err = trans.Commit().Error; err != nil {
				trans.Rollback()
			}
		}

	}()

	if trans.Error != nil {
		return 0, trans.Error
	}

	// 可以预见的错误，需要手动回滚
	if result := trans.Create(order); result.Error != nil {
		trans.Rollback()
		return 0, result.Error
	} else {
		rowAffected = result.RowsAffected
	}

	// 反序列化OrderData
	var deserializedData []models.OrderItem
	// 发生字符串拷贝
	err = json.Unmarshal([]byte(order.OrderData), &deserializedData)
	if err != nil {
		return 0, err
	}

	// 批量插入，自动开启事务
	if res := trans.Create(deserializedData); res.Error != nil {
		trans.Rollback()
		return 0, res.Error
	}

	return rowAffected, nil
}

// UpdateOrder 更新订单
//
// 可能需要手动更新orderitem表 考虑开启事务
func (o *OrderDAO) UpdateOrder(order *models.Order) (int64, error) {
	result := o.db.Save(order)
	return result.RowsAffected, result.Error
}

func (o *OrderDAO) GetOrderById(orderId int64) (*models.Order, error) {
	order := &models.Order{}
	result := o.db.Where("order_id = ?", orderId).First(order)
	return order, result.Error
}

func (o *OrderDAO) CreateOrderItem(orderItem *models.OrderItem) (int64, error) {
	result := o.db.Create(orderItem)
	return result.RowsAffected, result.Error
}
