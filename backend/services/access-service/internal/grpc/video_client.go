package grpc

import (
	"access-service/internal/types"
	pb "access-service/proto"
	"context"
	"errors"
	"google.golang.org/grpc"
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

func (c *Client) UserOwnsVideo(ctx context.Context, auth0ID string, videoID int64) (types.GrpcDBServiceResponse, error) {
	resp, err := c.DBService.CheckOwnership(ctx, &pb.CheckOwnershipRequest{
		Auth0Id: auth0ID,
		VideoId: videoID,
	})
	if err != nil {
		return types.GrpcDBServiceResponse{}, err
	}
	if resp.Message != "success" {
		return types.GrpcDBServiceResponse{}, errors.New("access denied: " + resp.Message)
	}
	return types.GrpcDBServiceResponse{
		Owned:     resp.Owned,
		Message:   resp.Message,
		ObjectKey: resp.ObjectKey,
		Bucket:    resp.Bucket,
	}, nil
}
