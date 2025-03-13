package grpc

import (
	pb "api-gateway/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn        *grpc.ClientConn
	UserService pb.UserServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:        conn,
		UserService: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
