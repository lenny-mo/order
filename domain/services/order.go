package services

import (
	"github.com/lenny-mo/order/domain/dao"
	"github.com/lenny-mo/order/domain/models"
)

type OrderServiceInterface interface {
	// 创建订单
	CreateOrder(order *models.Order) (int64, error)
	// 更新订单
	UpdateOrder(order *models.Order, oldversion int64) (int64, error)
	// 获取订单
	GetOrderById(orderId string) (*models.Order, error)
}

type OrderService struct {
	OrderDAO dao.OrderDAO
}

func NewOrderService(orderdao dao.OrderDAO) OrderServiceInterface {
	return &OrderService{
		OrderDAO: orderdao,
	}
}

func (o *OrderService) CreateOrder(order *models.Order) (int64, error) {
	return o.OrderDAO.CreateOrder(order)
}

func (o *OrderService) UpdateOrder(order *models.Order, oldversion int64) (int64, error) {
	return o.OrderDAO.UpdateOrder(order, oldversion)
}

func (o *OrderService) GetOrderById(orderId string) (*models.Order, error) {
	return o.OrderDAO.GetOrderById(orderId)
}
