package grpc

import (
	userPb "api-gateway/proto/user"
	"google.golang.org/grpc"
)

type Client struct {
	conn        *grpc.ClientConn
	UserService userPb.UserServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:        conn,
		UserService: userPb.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
