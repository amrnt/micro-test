package main

import (
	"github.com/amrnt/micro-test/api/pub/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

func main() {
	brkr := rabbitmq.NewBroker(
		broker.Addrs("amqp://rabbitmq:rabbitmq@localhost:5672"),
		rabbitmq.DurableExchange(),
		rabbitmq.Exchange("myexchange"),
	)

	if err := brkr.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := brkr.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.pub"),
		micro.Version("latest"),
		micro.Broker(brkr),
	)

	// Initialise service
	service.Init()

	// create new web service
	webService := web.NewService(
		web.Name("go.micro.api.pub"),
		web.Address(":3000"),
		web.Version("latest"),
		web.MicroService(service),
	)

	// initialise service
	if err := webService.Init(); err != nil {
		log.Fatal(err)
	}

	webService.HandleFunc("/", handler.Call(service))

	// Run service
	if err := webService.Run(); err != nil {
		log.Fatal(err)
	}
}
