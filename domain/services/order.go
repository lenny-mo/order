package services

import (
	"github.com/lenny-mo/order/domain/dao"
	"github.com/lenny-mo/order/domain/models"
)

type OrderServiceInterface interface {
	// 创建订单
	CreateOrder(order *models.Order) (int64, error)
	// 更新订单
	UpdateOrder(order *models.Order) (int64, error)
	// 获取订单
	GetOrderById(orderId int64) (*models.Order, error)
}

type OrderService struct {
	OrderDAO dao.OrderDAO
}

func NewOrderService(orderdao dao.OrderDAO) OrderServiceInterface {
	return &OrderService{
		OrderDAO: orderdao,
	}
}

// CreateOrder 创建订单
//
// # Order中的Orderdata 是OrderItem 列表的json序列化
//
// 需要反序列化成OrderItem列表
// 并且调用orderitem的插入方法
func (o *OrderService) CreateOrder(order *models.Order) (int64, error) {
	return o.OrderDAO.CreateOrder(order)
}

func (o *OrderService) UpdateOrder(order *models.Order) (int64, error) {
	return o.OrderDAO.UpdateOrder(order)
}

func (o *OrderService) GetOrderById(orderId int64) (*models.Order, error) {
	return o.OrderDAO.GetOrderById(orderId)
}
