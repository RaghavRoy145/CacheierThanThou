package client

import (
	"context"
	"fmt"
	"net"

	"github.com/raghavroy145/DistributedCaching/proto"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func New(endPoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endPoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func NewFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}
	resp, err := proto.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == proto.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key (%s)", key)
	}
	if resp.Status != proto.StatusOk {
		return nil, fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}
	return resp.Value, nil
}
func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}
	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", resp)
	if resp.Status != proto.StatusOk {
		return fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
