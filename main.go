package main

import (
	"fmt"
	"strconv"

	"github.com/lenny-mo/emall-utils/tracer"
	"github.com/lenny-mo/order/conf"
	"github.com/lenny-mo/order/domain/dao"
	"github.com/lenny-mo/order/domain/models"
	"github.com/lenny-mo/order/domain/services"
	"github.com/lenny-mo/order/handler"
	"github.com/lenny-mo/order/proto/order"
	"github.com/lenny-mo/order/utils"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func main() {
	// 1. 配置中心
	consulCof, err := conf.GetConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		fmt.Println(err)
		fmt.Println("获取配置失败")
		panic(err)
	}

	// 2. 注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	serviceName := "go.micro.service.order"
	// 3 链路追踪
	err = tracer.InitTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tracer.Closer.Close()
	opentracing.SetGlobalTracer(tracer.Tracer)
	// tracer, tracerio, err := utils.NewTracer("order-server", "127.0.0.1:6831")
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// defer tracerio.Close()
	// opentracing.SetGlobalTracer(tracer)

	// 4. 获取mysql配置
	mysqlConf := conf.GetMysqlFromConsul(consulCof, "mysql")

	// 5. 初始化数据库连接
	dsn := mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + strconv.FormatInt(mysqlConf.Port, 10) + ")/" + mysqlConf.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// 禁止复表的存在, 如果没有表则创建
	if !db.Migrator().HasTable(&models.Order{}) {
		db.Migrator().CreateTable(&models.Order{})
	}

	// 设置prometheus
	utils.PrometheusBoot(9091)

	// 创建服务
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8084"), // 服务监听地址
		// 使用consul注册中心
		micro.Registry(consulRegistry),
		// 添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// uber 漏桶 添加限流 每秒处理1000·个请求
		micro.WrapHandler(ratelimit.NewHandlerWrapper(conf.QPS)),
		// 添加prometheus
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	//
	service.Init()

	// 7. 创建service 和 handler 并且注册服务
	orderDAO := dao.NewOrderDAO(db)
	orderService := services.NewOrderService(*orderDAO.(*dao.OrderDAO))
	// 使用proto文件夹下的registry handler 方法注册
	err = order.RegisterOrderHandler(service.Server(), &handler.Order{
		Service: *orderService.(*services.OrderService),
	})
	if err != nil {
		panic(err)
	}

	// 8. 启动service
	if err = service.Run(); err != nil {
		fmt.Println(err)
		panic(err)
	}

}
