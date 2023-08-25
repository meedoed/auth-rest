package service

import (
	"context"
	"time"

	"github.com/meedoed/auth-rest/internal/domain"
	"github.com/meedoed/auth-rest/internal/repository"
	"github.com/meedoed/auth-rest/pkg/auth"
	"github.com/meedoed/auth-rest/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repo            repository.Users
	hasher          hash.PasswordHasher
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager, accessTTL, refreshTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *UsersService) GetTokens(ctx context.Context, guid string) (Tokens, error) {

	hex, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return Tokens{}, err
	}

	if _, err := s.repo.GetById(ctx, hex); err != nil {
		return Tokens{}, nil
	}

	tokens, err := s.createSession(ctx, guid)
	if err != nil {
		return Tokens{}, err
	}

	return tokens, nil
}

func (s *UsersService) RefreshTokens(ctx context.Context, accessToken, refreshToken string) (Tokens, error) {
	guid, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return Tokens{}, err
	}

	hex, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return Tokens{}, err
	}

	session, err := s.repo.GetById(ctx, hex)
	if err != nil {
		return Tokens{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), []byte(refreshToken)); err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, guid)
}

func (s *UsersService) createSession(ctx context.Context, guid string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(guid, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	hashedToken, err := hash.HashToken(res.RefreshToken)
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: hashedToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	hex, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return res, err
	}

	err = s.repo.SetSession(ctx, hex, session)

	return res, err
}
