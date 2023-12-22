package handler

import (
	"context"
	"errors"
	"time"

	m "github.com/lenny-mo/emall-utils/metrics"
	"github.com/lenny-mo/order/domain/models"
	"github.com/lenny-mo/order/domain/services"
	"github.com/lenny-mo/order/proto/order"
	"github.com/lenny-mo/order/utils"
)

// 需要实现的接口
//
//	type OrderHandler interface {
//		// 插入操作涉及到幂等性，需要生成全局唯一的订单ID
//		InsertOrder(context.Context, *InserRequest, *InserResponse) error
//		GetOrder(context.Context, *GetRequest, *GetResponse) error
//		// 更新操作涉及到乐观锁做并发控制，所以需要传入版本号
//		UpdateOrder(context.Context, *UpdateRequest, *UpdateResponse) error
//		// 用户在创建订单的时候需要先调用此方法生成订单号
//		GenerateUUID(context.Context, *Empty, *GenerateUUIDResponse) error
//	}
type Order struct {
	Service services.OrderService
}

// 用于指定prometheus监控label
const (
	SERVICE = "order"
	VERSION = "v1.0.0"
	OS      = "centOS"
)

func (o *Order) InsertOrder(ctx context.Context, req *order.InserRequest, res *order.InserResponse) error {
	// prometheus 请求数+1
	m.CounterRequestProcess(SERVICE, VERSION, OS)
	// 记录支付处理开始时间
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		m.RecordPaymentResponseTime(SERVICE, VERSION, OS, duration) // Histogram 指标
		m.TaskExecutionTime(SERVICE, VERSION, OS, duration)         // Summary 指标
	}()

	order := &models.Order{
		UserId:       req.OrderData.UserId,
		OrderData:    req.OrderData.OrderData,
		OrderId:      req.OrderData.OrderId,
		Status:       int8(req.OrderData.Status),
		OrderVersion: req.OrderData.OrderVersion,
	}
	rowAffected, err := o.Service.CreateOrder(order)
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("insert order failed, row affected is 0")
	}

	res.RowsAffected = int32(rowAffected)

	return nil
}

func (o *Order) GetOrder(ctx context.Context, req *order.GetRequest, res *order.GetResponse) error {
	// prometheus 请求数+1
	m.CounterRequestProcess(SERVICE, VERSION, OS)
	// 记录支付处理开始时间
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		m.RecordPaymentResponseTime(SERVICE, VERSION, OS, duration) // Histogram 指标
		m.TaskExecutionTime(SERVICE, VERSION, OS, duration)         // Summary 指标
	}()

	orderdata, err := o.Service.GetOrderById(req.OrderId)
	if err != nil {
		return err
	}

	res.OrderData = &order.OrderInfo{
		OrderId:      orderdata.OrderId,
		UserId:       orderdata.UserId,
		OrderVersion: orderdata.OrderVersion,
		OrderData:    orderdata.OrderData,
		Status:       order.OrderStatus(orderdata.Status),
	}

	return nil
}

func (o *Order) UpdateOrder(ctx context.Context, req *order.UpdateRequest, res *order.UpdateResponse) error {
	// prometheus 请求数+1
	m.CounterRequestProcess(SERVICE, VERSION, OS)
	// 记录支付处理开始时间
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		m.RecordPaymentResponseTime(SERVICE, VERSION, OS, duration) // Histogram 指标
		m.TaskExecutionTime(SERVICE, VERSION, OS, duration)         // Summary 指标
	}()

	order := &models.Order{
		UserId:       req.OrderData.UserId,
		OrderData:    req.OrderData.OrderData,
		OrderId:      req.OrderData.OrderId,
		Status:       int8(req.OrderData.Status),
		OrderVersion: req.OrderData.OrderVersion,
	}
	rowAffected, err := o.Service.UpdateOrder(order, req.Oldversion)
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		res.RowsAffected = 0
		return errors.New("update order failed, row affected is 0")
	}
	res.RowsAffected = int32(rowAffected)
	return nil
}

func (o *Order) GenerateUUID(ctx context.Context, req *order.Empty, res *order.GenerateUUIDResponse) error {
	res.Uuid = utils.UUID()
	return nil
}
