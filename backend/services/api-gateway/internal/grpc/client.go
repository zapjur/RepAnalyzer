package grpc

import (
	dbPb "api-gateway/proto/db"
	"google.golang.org/grpc"
)

type Client struct {
	conn      *grpc.ClientConn
	DBService dbPb.DBServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      conn,
		DBService: dbPb.NewDBServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
