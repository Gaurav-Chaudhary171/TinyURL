package service

import (
	"context"
	"fmt"
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"
)

type ShortenServer struct {
	proto.UnimplementedTinyURLServiceServer
}

func (s *ShortenServer) ShortenURL(ctx context.Context, req *proto.ShortenURLRequest) (*proto.ShortenURLResponse, error) {
	if req.Url == "" {
		return nil, fmt.Errorf("url is required")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	shortUrl := ShortenURL(req.Url)
	shortUrl = "https://" + shortUrl

	url := &model.GeneratedUrl{
		OriginalUrl: req.Url,
		Username:    req.Username,
		TinyUrl:     shortUrl,
		IsActive:    true,
	}

	if err := config.DB.Create(url).Error; err != nil {
		return nil, fmt.Errorf("error creating shortened url: %v", err)
	}

	return &proto.ShortenURLResponse{
		Status:      "success",
		Shortenurl:  shortUrl,
		Originalurl: req.Url,
	}, nil
}
