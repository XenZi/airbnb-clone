package services

import (
	"auth-service/domains"
	"auth-service/repository"
	errors2 "errors"
)

type UserService struct {
	userRepository  *repository.UserRepository
	passwordService *PasswordService
	jwtService      *JwtService
}

func NewUserService(userRepo *repository.UserRepository, passwordService *PasswordService, jwtService *JwtService) *UserService {
	return &UserService{
		userRepository:  userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}
func (u *UserService) CreateUser(registerUser domains.RegisterUser) (*domains.UserDTO, error) {
	user := domains.User{
		Email:    registerUser.Email,
		Password: registerUser.Password,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	newUser, err := u.userRepository.SaveUser(user)

	if err != nil {
		return nil, err
	}
	id, err := newUser.ID.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &domains.UserDTO{
		ID:       string(id),
		Email:    registerUser.Email,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}, nil
}

func (u *UserService) LoginUser(loginData domains.LoginUser) (*string, error) {
	user, err := u.userRepository.FindUserByEmail(loginData.Email)
	if err != nil {
		return nil, err
	}
	isSamePassword := u.passwordService.CheckPasswordHash(loginData.Password, user.Password)
	if !isSamePassword {
		return nil, errors2.New("Bad credentials")
	}
	jwtToken, err := u.jwtService.CreateKey(user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}
