package services

import (
	"auth-service/client"
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/repository"
	"auth-service/utils"
	"context"
	"log"
)

type UserService struct {
	userRepository    *repository.UserRepository
	passwordService   *PasswordService
	jwtService        *JwtService
	validator         *utils.Validator
	encryptionService *EncryptionService
	mailClient        client.MailClientInterface
	userClient 		  *client.UserClient
}

func NewUserService(userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JwtService,
	validator *utils.Validator,
	encryptionService *EncryptionService,
	mailClient client.MailClientInterface,
	userClient *client.UserClient) *UserService {
	return &UserService{
		userRepository:    userRepo,
		passwordService:   passwordService,
		jwtService:        jwtService,
		validator:         validator,
		encryptionService: encryptionService,
		mailClient:        mailClient,
		userClient: userClient,
	}
}
func (u *UserService) CreateUser(ctx context.Context, registerUser domains.RegisterUser) (*domains.UserDTO, *errors.ErrorStruct) {
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
	log.Println(registerUser.Age)
	go func() {
		u.mailClient.SendAccountConfirmationEmail(registerUser.Email, token)
	}()
	u.userClient.SendCreatedUser(ctx, newUser.ID.Hex(), registerUser)
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
	jwtToken, foundError := u.jwtService.CreateKey(user.Email, user.Role, user.ID.Hex())
	if foundError != nil {
		return nil, errors.NewError(
			foundError.Error(),
			500)
	}
	return &domains.SuccessfullyLoggedUser{
		Token: *jwtToken,
		User: domains.UserDTO{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (u UserService) ConfirmUserAccount(token string) (*domains.UserDTO, *errors.ErrorStruct) {
	userID, err := u.encryptionService.ValidateToken(token)
	if err != nil {
		log.Println(err.GetErrorMessage())
		return nil, err
	}
	updatedUser, err := u.userRepository.UpdateUserConfirmation(userID)
	if err != nil {
		log.Println(err.GetErrorMessage())
		return nil, err
	}
	return &domains.UserDTO{
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		ID:        userID,
		Role:      updatedUser.Role,
		Confirmed: updatedUser.Confirmed,
	}, nil
}

func (u UserService) RequestResetPassword(email string) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	if email == "" {
		return nil, errors.NewError("Email is empty or incorrect", 400)
	}
	user, err := u.userRepository.FindUserByEmail(email)
	if err != nil {
		log.Println(err)
		return nil, errors.NewError(err.GetErrorMessage(), err.GetErrorStatus())
	}
	token, err := u.encryptionService.GenerateToken(user.ID.Hex())
	if err != nil {
		log.Println(err)
		return nil, errors.NewError(err.GetErrorMessage(), err.GetErrorStatus())
	}
	go func() {
		u.mailClient.SendRequestResetPassword(email, token)
	}()
	return &domains.BaseMessageResponse{
		Message: "Email has been sent",
	}, nil
}

func (u UserService) ResetPassword(requestData domains.ResetPassword, token string) (*domains.UserDTO, *errors.ErrorStruct) {
	u.validator.ValidatePassword(requestData.Password)
	validatorErrors := u.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	userID, err := u.encryptionService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if requestData.Password != requestData.ConfirmedPassword {
		return nil, errors.NewError("Password doesn't match", 400)
	}
	if u.passwordService.CheckPasswordExistanceInBlacklist(requestData.Password) != false {
		return nil, errors.NewError("Choose better password that is more secure!", 400)
	}
	hashedPassword, hashError := u.passwordService.HashPassword(requestData.Password)
	if hashError != nil {
		return nil, err
	}
	user, err := u.userRepository.UpdateUserPassword(userID, hashedPassword)
	return &domains.UserDTO{
		Username:  user.Username,
		Email:     user.Email,
		ID:        userID,
		Role:      user.Role,
		Confirmed: user.Confirmed,
	}, nil
}

func (u UserService) ChangePassword(data domains.ChangePassword, userID string) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	if data.ConfirmedPassword != data.Password {
		return nil, errors.NewError("New password doesn't match with each other", 400)
	}
	if u.passwordService.CheckPasswordExistanceInBlacklist(data.Password) {
		return nil, errors.NewError("Choose better password", 400)
	}
	user, err := u.userRepository.FindUserById(userID)

	if err != nil {
		return nil, err
	}

	if u.passwordService.CheckPasswordHash(data.Password, user.Password) == true {
		return nil, errors.NewError("Old password doesn't match", 400)
	}

	hashedPassword, hashError := u.passwordService.HashPassword(data.Password)
	if hashError != nil {
		return nil, errors.NewError(hashError.Error(), 500)
	}
	_, err = u.userRepository.UpdateUserPassword(userID, hashedPassword)
	if err != nil {
		return nil, err
	}
	response := domains.BaseMessageResponse{
		Message: "You have updated your password",
	}
	return &response, nil
}
