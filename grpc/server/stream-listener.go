package server

import (
	"context"
	"errors"

	"github.com/ajainc/chain/grpc/gen"
)

type StreamListener struct {
	c          context.Context
	cancelFunc context.CancelFunc
	authed     gen.Authed
	mailBox    chan []byte
	sink       gen.StreamService_ConnectServer
}

func newStreamListener(c context.Context, authed *gen.Authed, mailBoxSize int, sink gen.StreamService_ConnectServer) (*StreamListener, error) {
	if authed == nil {
		return nil, errors.New("no auth")
	}

	c, cancel := context.WithCancel(c)
	return &StreamListener{
		c:          c,
		cancelFunc: cancel,
		authed:     *authed,
		mailBox:    make(chan []byte, mailBoxSize),
		sink:       sink,
	}, nil
}

func (rcvr *StreamListener) Listen() {
	for {
		select {
		case message := <-rcvr.mailBox:
			println(message)
			//TODO
			//sink.Send()
		case <-rcvr.c.Done():
			return
		}
	}
}

func (rcvr *StreamListener) Send(data []byte) {
	rcvr.mailBox <- data
}

func (rcvr *StreamListener) Close() {
	rcvr.cancelFunc()
}
