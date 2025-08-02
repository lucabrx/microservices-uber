package auth

import (
	"context"

	pb "github.com/lukabrx/uber-clone/api/proto/auth/v1"
)

type GrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	service *Service
}

func NewGrpcHandler(service *Service) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) AuthenticateWithGoogle(ctx context.Context, req *pb.AuthenticateWithGoogleRequest) (*pb.AuthenticateWithGoogleResponse, error) {
	token, err := h.service.AuthenticateWithGoogle(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.AuthenticateWithGoogleResponse{Token: token}, nil
}

func (h *GrpcHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	payload, err := h.service.VerifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResponse{UserId: payload.UserID}, nil
}
