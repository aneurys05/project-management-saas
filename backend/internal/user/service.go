package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aneurys05/project-management-saas/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
	jwtSecret  string
}

func NewService(repository *Repository, jwtSecret string) *Service {
	return &Service{
		repository: repository,
		jwtSecret:  jwtSecret,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

func (s *Service) Register(
	ctx context.Context,
	input RegisterInput,
) (*User, error) {

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)

	if input.Email == "" {
		return nil, errors.New("email is required")
	}

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if len(input.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return s.repository.Create(
		ctx,
		input.Email,
		string(passwordHash),
		input.Name,
	)
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *Service) Login(
	ctx context.Context,
	input LoginInput,
) (string, *User, error) {

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.repository.FindByEmail(ctx, input.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := auth.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}

	return token, user, nil
}
