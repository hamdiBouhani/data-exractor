package grpc

import (
	"context"

	"sm-connector-be/pb"
)

type Channel struct {
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}
type JwtUser struct {
	Sub   string `protobuf:"bytes,1,opt,name=sub,proto3" json:"sub,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}
type GetMessageRequest struct {
	Content   string   `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	MessageId string   `protobuf:"bytes,2,opt,name=messageId,proto3" json:"messageId,omitempty"`
	Channel   *Channel `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel,omitempty"`
	SendTo    *JwtUser `protobuf:"bytes,4,opt,name=sendTo,proto3" json:"sendTo,omitempty"`
}

// Ping pong service
func (s *Server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{
		Pong: "pong",
	}, nil
}

// method used to emit the message content
func (s *Server) GetMessage(ctx context.Context, in *pb.GetMessageRequest) (*pb.GetMessageResponse, error) {

	return &pb.GetMessageResponse{}, nil
}
