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
	accessToken, refreshToken, user, err := h.service.AuthenticateWithGoogle(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateWithGoogleResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &pb.User{
			Id:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (h *GrpcHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	payload, err := h.service.VerifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResponse{UserId: payload.UserID}, nil
}

func (h *GrpcHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	newAccessToken, newRefreshToken, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (h *GrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	pbUser := &pb.User{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return &pb.GetUserResponse{User: pbUser}, nil
}
