package main

import (
	"fmt"

	"github.com/amrnt/micro-test/srv/sub/handler"
	"github.com/amrnt/micro-test/srv/sub/subscriber"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/broker/rabbitmq"

	sub "github.com/amrnt/micro-test/srv/sub/proto/sub"
	gclient "github.com/micro/go-micro/client/grpc"
	gserver "github.com/micro/go-micro/server/grpc"
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

	var service micro.Service

	// New Service

	// Please comment out one of the service
	service = grpc.NewService(
		micro.Name("go.micro.srv.sub"),
		micro.Version("latest"),
		micro.Broker(brkr),
	)

	// Same as grpc.NewService above
	service = micro.NewService(
		micro.Name("go.micro.srv.sub"),
		micro.Version("latest"),
		micro.Broker(brkr),
		micro.Client(gclient.NewClient()),
		micro.Server(gserver.NewServer()),
	)

	// Comment out this to see that
	// `RegisterSubscriber` is not working anymore
	service = micro.NewService(
		micro.Name("go.micro.srv.sub"),
		micro.Version("latest"),
		micro.Broker(brkr),
	)

	// Initialise service
	service.Init()

	// Register Handler
	sub.RegisterSubHandler(service.Server(), new(handler.Sub))

	// ---
	// Subscribing in 2 ways
	// ---

	// 1. with micro.RegisterSubscriber
	// Works only with micro.NewService
	// When it comes to respect the context and the parameters
	// Passed thru server.SubscriberContext

	// Need to pass more options to the sub
	brkrSub := broker.NewSubscribeOptions(
		broker.Queue("go.micro.srv.sub.topic.1.default"),
		rabbitmq.DurableQueue(),
		broker.DisableAutoAck(),
		rabbitmq.QueueArguments(map[string]interface{}{"x-queue-type": "quorum"}),
		rabbitmq.RequeueOnError(),
		// rabbitmq.AckOnSuccess(),
	)

	// Register Struct as Subscriber 1
	// Set Context from broker.NewSubscribeOptions
	// Should be Durable
	// With arguments
	micro.RegisterSubscriber(
		"go.micro.srv.sub.topic.1",
		service.Server(),
		new(subscriber.Sub),
		server.SubscriberContext(brkrSub.Context),
		server.SubscriberQueue("go.micro.srv.sub.topic.1.default"),
		server.DisableAutoAck(),
	)

	// Register Struct as Subscriber 2
	// This will be auto-delete queue
	micro.RegisterSubscriber(
		"go.micro.srv.sub.topic.2",
		service.Server(),
		new(subscriber.Sub),
		server.SubscriberQueue("go.micro.srv.sub.topic.2.default"),
		server.DisableAutoAck(),
	)

	// Register Function as Subscriber
	// micro.RegisterSubscriber("go.micro.srv.sub.topic.2", service.Server(), subscriber.Handler)

	// 2. with service.Options().Broker.Subscribe()
	// Workis fine with the micro.NewService and grpc.NewService

	sx, err := service.Options().Broker.Subscribe(
		"go.micro.srv.sub.topic.3",
		func(p broker.Event) error {
			fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
			return nil
			// return fmt.Errorf("err")
		},
		broker.Queue("go.micro.srv.sub.topic.3.default"),
		rabbitmq.DurableQueue(),
		broker.DisableAutoAck(),
		rabbitmq.QueueArguments(map[string]interface{}{"x-queue-type": "quorum"}),
		rabbitmq.RequeueOnError(),
	)
	if err != nil {
		fmt.Println("service.Options().Broker.Subscribe", err)
	}
	defer sx.Unsubscribe()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
