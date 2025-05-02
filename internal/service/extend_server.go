package service

import (
	"context"
	"fmt"
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"
)

type ExtendServer struct {
	proto.UnimplementedTinyURLServiceServer
}

func (s *ExtendServer) ExtendURL(ctx context.Context, req *proto.ExtendURLRequest) (*proto.ExtendURLResponse, error) {
	if req.Url == "" {
		return nil, fmt.Errorf("url is required")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	var url model.GeneratedUrl
	result := config.DB.Where("tiny_url = ? AND username = ? AND is_active = ?", req.Url, req.Username, true).First(&url)
	if result.Error != nil {
		return nil, fmt.Errorf("error extending url: %v", result.Error)
	}

	return &proto.ExtendURLResponse{
		Status:      "success",
		Originalurl: url.OriginalUrl,
		Extenedurl:  req.Url,
	}, nil
}
