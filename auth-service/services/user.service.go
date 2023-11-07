package services

import (
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/repository"
	"auth-service/utils"
)

type UserService struct {
	userRepository  *repository.UserRepository
	passwordService *PasswordService
	jwtService      *JwtService
	validator       *utils.Validator
}

func NewUserService(userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JwtService,
	validator *utils.Validator) *UserService {
	return &UserService{
		userRepository:  userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		validator:       validator,
	}
}
func (u *UserService) CreateUser(registerUser domains.RegisterUser) (*domains.UserDTO, *errors.ErrorStruct) {
	u.validator.ValidateRegisterUser(&registerUser)
	validatorErrors := u.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	user := domains.User{
		Email:    registerUser.Email,
		Password: registerUser.Password,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	user.Password = hashedPassword
	newUser, foundErr := u.userRepository.SaveUser(user)
	if foundErr != nil {
		return nil, foundErr
	}
	id, err := newUser.ID.MarshalJSON()
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return &domains.UserDTO{
		ID:       string(id),
		Email:    registerUser.Email,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}, nil
}

func (u *UserService) LoginUser(loginData domains.LoginUser) (*string, *errors.ErrorStruct) {
	user, err := u.userRepository.FindUserByEmail(loginData.Email)
	if err != nil {
		return nil, err
	}
	isSamePassword := u.passwordService.CheckPasswordHash(loginData.Password, user.Password)
	if !isSamePassword {
		return nil, errors.NewError("Bad credentials", 401)
	}
	jwtToken, foundError := u.jwtService.CreateKey(user.Email, user.Role)
	if foundError != nil {
		return nil, errors.NewError(
			foundError.Error(),
			500)
	}
	return jwtToken, nil
}
