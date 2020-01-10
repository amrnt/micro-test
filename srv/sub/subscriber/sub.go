package subscriber

import (
	"context"

	"github.com/micro/go-micro/util/log"

	sub "github.com/amrnt/micro-test/srv/sub/proto/sub"
)

type Sub struct{}

func (e *Sub) Handle(ctx context.Context, msg *sub.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	// return fmt.Errorf("err")
	return nil
}

func Handler(ctx context.Context, msg *sub.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
