package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	sub "github.com/amrnt/micro-test/srv/sub/proto/sub"
)

type Sub struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Sub) Call(ctx context.Context, req *sub.Request, rsp *sub.Response) error {
	log.Log("Received Sub.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Sub) Stream(ctx context.Context, req *sub.StreamingRequest, stream sub.Sub_StreamStream) error {
	log.Logf("Received Sub.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&sub.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Sub) PingPong(ctx context.Context, stream sub.Sub_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&sub.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
