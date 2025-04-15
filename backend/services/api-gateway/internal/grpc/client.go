package grpc

import (
	orPb "api-gateway/proto/analysis"
	dbPb "api-gateway/proto/db"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	conn                *grpc.ClientConn
	DBService           dbPb.DBServiceClient
	OrchestratorService orPb.OrchestratorClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := ConnectGRPC(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:                conn,
		DBService:           dbPb.NewDBServiceClient(conn),
		OrchestratorService: orPb.NewOrchestratorClient(conn),
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
