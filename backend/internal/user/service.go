package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
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
