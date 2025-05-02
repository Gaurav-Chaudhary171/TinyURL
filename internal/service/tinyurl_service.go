package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"TinyURL_Refactored/proto"

	"gorm.io/gorm"
)

type TinyURLService struct {
	proto.UnimplementedTinyURLServiceServer
	db *gorm.DB
}

func NewTinyURLService() *TinyURLService {
	return &TinyURLService{
		db: config.DB,
	}
}

// ShortenURL implements the URL shortening functionality
func (s *TinyURLService) ShortenURL(ctx context.Context, req *proto.ShortenURLRequest) (*proto.ShortenURLResponse, error) {
	if req.Url == "" {
		return nil, fmt.Errorf("url is required")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	shortUrl := generateShortURL(req.Url)
	shortUrl = "https://" + shortUrl

	url := &model.GeneratedUrl{
		OriginalUrl: req.Url,
		Username:    req.Username,
		TinyUrl:     shortUrl,
		IsActive:    true,
	}

	if err := s.db.Create(url).Error; err != nil {
		return nil, fmt.Errorf("error creating shortened url: %v", err)
	}

	return &proto.ShortenURLResponse{
		Status:      "success",
		Shortenurl:  shortUrl,
		Originalurl: req.Url,
	}, nil
}

// ExtendURL implements the URL extension functionality
func (s *TinyURLService) ExtendURL(ctx context.Context, req *proto.ExtendURLRequest) (*proto.ExtendURLResponse, error) {
	if req.Url == "" {
		return nil, fmt.Errorf("url is required")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}

	var url model.GeneratedUrl
	result := s.db.Where("tiny_url = ? AND username = ? AND is_active = ?", req.Url, req.Username, true).First(&url)
	if result.Error != nil {
		return nil, fmt.Errorf("error extending url: %v", result.Error)
	}

	return &proto.ExtendURLResponse{
		Status:      "success",
		Originalurl: url.OriginalUrl,
		Extenedurl:  req.Url,
	}, nil
}

// RegisterUser implements the user registration functionality
func (s *TinyURLService) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if req.FirstName == "" {
		return nil, fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return nil, fmt.Errorf("last name is required")
	}

	var existingUser model.Users
	result := s.db.Where("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, fmt.Errorf("username already exists")
	}

	user := model.Users{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("error registering user: %v", err)
	}

	return &proto.RegisterUserResponse{
		Status:    "success",
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// Login implements the user login functionality
func (s *TinyURLService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	var user model.Users
	result := s.db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found: %v", result.Error)
	}

	return &proto.LoginResponse{
		Status: "success",
		User: &proto.RegisterUserRequest{
			Username: user.Username,
		},
	}, nil
}

// generateShortURL generates a shortened URL from the original URL
func generateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8] // Return first 8 characters
}
