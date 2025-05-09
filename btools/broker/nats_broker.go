package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/seemyown/backend-toolkit/btools/logging"
	"github.com/seemyown/nats-rpc-go/natsrpc"
)

type Publisher interface {
	Publish(topic string, data map[string]interface{}) error
	Request(ctx context.Context, topic string, hdr nats.Header, in, out interface{}) error
	Close()
}

type natsBroker struct {
	nc  *nats.Conn
	log *logging.Logger
}

func NewNatsBroker(natsConn *nats.Conn, logger *logging.Logger) Publisher {
	return &natsBroker{nc: natsConn, log: logger}
}

func (n *natsBroker) Publish(topic string, data map[string]interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		n.log.Error(err, "error marshal data")
		return err
	}
	err = n.nc.Publish(topic, msg)
	if err != nil {
		n.log.Error(err, "error send message %s in %s", data, topic)
		return err
	}
	n.log.Debug("message %s successful sent in %s", data, topic)
	return nil
}

func (n *natsBroker) Request(ctx context.Context, topic string, hdr nats.Header, in, out interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	msg := &nats.Msg{
		Subject: topic,
		Data:    data,
		Header:  hdr,
	}

	resp, err := n.nc.RequestMsgWithContext(ctx, msg)
	if err != nil {
		return err
	}

	var rpcErr natsrpc.RPCError
	if err := json.Unmarshal(resp.Data, &rpcErr); err == nil && rpcErr.Code != 0 {
		return &rpcErr
	}

	if err := json.Unmarshal(resp.Data, out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

func (n *natsBroker) Close() {
	n.nc.Close()
}
