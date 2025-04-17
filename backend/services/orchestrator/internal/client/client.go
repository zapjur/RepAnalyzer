package client

import (
	"google.golang.org/grpc"
	pb "orchestrator/proto/db"
	"time"
)

type Client struct {
	conn      *grpc.ClientConn
	DBService pb.DBServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := ConnectGRPC(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      conn,
		DBService: pb.NewDBServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func ConnectGRPC(address string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	for i := 0; i < 5; i++ { // Retry 5 times
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err == nil {
			return conn, nil
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	return nil, err
}
