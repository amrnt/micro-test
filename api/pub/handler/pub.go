package handler

import (
	"context"
	"net/http"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/streadway/amqp"

	sub "github.com/amrnt/micro-test/srv/sub/proto/sub"
)

// Pub ...
type Pub struct{}

// Call ...
func Call(service micro.Service, p micro.Publisher, p2 micro.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error

		// Publishing 2 ways

		// ---
		// 1. Using Broker.Publish
		// to go.micro.srv.sub.topic.3
		// ---

		msg := &broker.Message{
			Header: map[string]string{"id": "1"},
			Body:   []byte("hello"),
		}

		err = service.Options().Broker.Publish("go.micro.srv.sub.topic.3", msg, rabbitmq.DeliveryMode(amqp.Persistent))
		if err != nil {
			log.Errorf("go.micro.srv.sub.topic.3: error publishing: %s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Something went wrong publishing go.micro.srv.sub.topic.3"))
			return
		}

		// ---
		// 2. Using publisher.Publish
		// to go.micro.srv.sub.topic.1
		// ---

		// Try to make a new context to pass some pub options
		pubOpt := NewPublsiherOptions(
			rabbitmq.DeliveryMode(amqp.Persistent),
		)

		// setting or passing pubOpt Context is not working
		// BUT event is being published if you comment out the options
		// from below
		err = p.Publish(
			pubOpt.Context, // context.Background()
			&sub.Message{Say: "Hello world!"},
			func(options *client.PublishOptions) {
				// options.Exchange = "myexchange"
				// options.Context = pubOpt.Context
			},
		)
		if err != nil {
			log.Errorf("go.micro.srv.sub.topic.1: error publishing: %s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Something went wrong publishing go.micro.srv.sub.topic.1"))
			return
		}

		// Same thing but publish to topic 2
		// Subscriber doesnt set any new context
		err = p2.Publish(
			pubOpt.Context, // context.Background()
			&sub.Message{Say: "Hello world!"},
			func(options *client.PublishOptions) {
				// options.Exchange = "myexchange"
				// options.Context = pubOpt.Context
			},
		)
		if err != nil {
			log.Errorf("go.micro.srv.sub.topic.2: error publishing: %s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Something went wrong publishing go.micro.srv.sub.topic.2"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}
}

// NewPublsiherOptions ...
func NewPublsiherOptions(opts ...broker.PublishOption) broker.PublishOptions {
	opt := broker.PublishOptions{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
