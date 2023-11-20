package services

import (
	"auth-service/client"
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/repository"
	"auth-service/utils"
	"log"
)

type UserService struct {
	userRepository    *repository.UserRepository
	passwordService   *PasswordService
	jwtService        *JwtService
	validator         *utils.Validator
	encryptionService *EncryptionService
	mailClient        client.MailClientInterface
}

func NewUserService(userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JwtService,
	validator *utils.Validator,
	encryptionService *EncryptionService,
	mailClient client.MailClientInterface) *UserService {
	return &UserService{
		userRepository:    userRepo,
		passwordService:   passwordService,
		jwtService:        jwtService,
		validator:         validator,
		encryptionService: encryptionService,
		mailClient:        mailClient,
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
	if u.passwordService.CheckPasswordExistanceInBlacklist(registerUser.Password) != false {
		return nil, errors.NewError("Choose better password that is more secure!", 400)
	}
	user := domains.User{
		Email:     registerUser.Email,
		Password:  registerUser.Password,
		Username:  registerUser.Username,
		Role:      registerUser.Role,
		Confirmed: false,
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
	token, encError := u.encryptionService.GenerateToken(string(id[1 : len(id)-1]))
	if encError != nil {
		return nil, encError
	}

	go func() {
		u.mailClient.SendAccountConfirmationEmail(registerUser.Email, token)
	}()
	return &domains.UserDTO{
		ID:       string(id[1 : len(id)-1]),
		Email:    registerUser.Email,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}, nil
}

func (u *UserService) LoginUser(loginData domains.LoginUser) (*domains.SuccessfullyLoggedUser, *errors.ErrorStruct) {
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
	id, _ := user.ID.MarshalJSON()
	return &domains.SuccessfullyLoggedUser{
		Token: *jwtToken,
		User: domains.UserDTO{
			ID:       string(id),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}


func (u UserService) ConfirmUserAccount(token string) (*domains.BaseHttpResponse, *errors.ErrorStruct) {
	log.Println(token)
	userID, err := u.encryptionService.ValidateToken(token)
	if err != nil {
		log.Println(err.GetErrorMessage())
		return nil, err
	}
	log.Println(userID)
	return nil, nil
}
