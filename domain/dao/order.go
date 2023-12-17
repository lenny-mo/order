package dao

import (
	"errors"
	"fmt"

	"github.com/lenny-mo/order/domain/models"
	"gorm.io/gorm"
)

type OrderDAOInterface interface {
	CreateOrder(order *models.Order) (int64, error)
	UpdateOrder(order *models.Order, oldversion int64) (int64, error)
	GetOrderById(orderId string) (*models.Order, error)
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

// CreateOrder 创建订单
func (o *OrderDAO) CreateOrder(order *models.Order) (rowAffected int64, err error) {

	result := o.db.Create(order)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.RowsAffected, err
	}

	return result.RowsAffected, nil
}

// UpdateOrder 更新订单
func (o *OrderDAO) UpdateOrder(order *models.Order, oldversion int64) (int64, error) {
	// 1. 获取当前orderid的记录，检查version 是否被修改, 如果被修改，则返回0, err 表示没有更新，并且告诉上游发现数据库更新冲突
	oldData, err := o.GetOrderById(order.OrderId)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	if oldData.OrderVersion != oldversion {
		fmt.Println(err)
		return 0, errors.New("update order failed, found confliction while on order version while update")
	}
	// 2. 正常更新
	result := o.db.Model(&models.Order{}).Where("order_id= ?", order.OrderId).Updates(order)
	return result.RowsAffected, result.Error
}

func (o *OrderDAO) GetOrderById(orderId string) (*models.Order, error) {
	order := &models.Order{}
	result := o.db.Where("order_id = ?", orderId).First(order)
	return order, result.Error
}

func (o *OrderDAO) CreateOrderItem(orderItem *models.OrderItem) (int64, error) {
	result := o.db.Create(orderItem)
	return result.RowsAffected, result.Error
}
