package services

import (
	"auth-service/client"
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/repository"
	"auth-service/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepository    *repository.UserRepository
	passwordService   *PasswordService
	jwtService        *JwtService
	validator         *utils.Validator
	encryptionService *EncryptionService
	mailClient        client.MailClientInterface
	userClient 		  *client.UserClient
	notificationClient *client.NotificationClient
}

func NewUserService(userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JwtService,
	validator *utils.Validator,
	encryptionService *EncryptionService,
	mailClient client.MailClientInterface,
	userClient *client.UserClient,
	notificationClient *client.NotificationClient) *UserService {
	return &UserService{
		userRepository:    userRepo,
		passwordService:   passwordService,
		jwtService:        jwtService,
		validator:         validator,
		encryptionService: encryptionService,
		mailClient:        mailClient,
		userClient: userClient,
		notificationClient: notificationClient,
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
	errFromUserService := u.userClient.SendCreatedUser(ctx, newUser.ID.Hex(), registerUser)
	if errFromUserService != nil {
		_, err := u.userRepository.DeleteUserById(newUser.ID.Hex())
		if err != nil {
			return nil, err
		}
		return nil, errFromUserService
	}
	go func() {
		u.mailClient.SendAccountConfirmationEmail(registerUser.Email, token)
	}()
	errFromNotificationService := u.notificationClient.CreateNewUserStructNotification(ctx, newUser.ID.Hex())
	if errFromNotificationService != nil {
		log.Println("ERR FROM NOTIFICATION", errFromNotificationService)
	}
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
	if u.passwordService.CheckPasswordHash(data.OldPassword, user.Password) == false {
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

func (u UserService) UpdateCredentials(ctx context.Context, id string, updatedData domains.User) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	if (updatedData.Email == "" || updatedData.Username == "") {
		return nil, errors.NewError("Email or username are empty", 400)
	} 
	foundUser, _ := u.userRepository.FindUserById(id)
	if (foundUser.Username == updatedData.Username && foundUser.Email == updatedData.Email) {
		return &domains.BaseMessageResponse{
			Message: "You have successfully updated your credentials", 
		}, nil
	}
	if (foundUser.Email != updatedData.Email) {
		_, err := u.userRepository.FindUserByEmail(updatedData.Email)
		if err == nil {
			return  nil, errors.NewError("User with same email already exists", 400)
		}
	}
	if (foundUser.Username != updatedData.Username) {
		_, err := u.userRepository.FindUserByUsername(updatedData.Username)
		if err == nil {
			return nil, errors.NewError("User with same username already exists", 400)
		}
	}
	if u.passwordService.CheckPasswordHash(updatedData.Password, foundUser.Password) == false {
		return nil, errors.NewError("Password doesn't match", 400)
	}
	objectID, newError := primitive.ObjectIDFromHex(id)
	if newError != nil {
		return nil, errors.NewError(newError.Error(),500)
	}
	updatedData.ID = objectID
	errFromCredentialsUpdate := u.userClient.SendUpdateCredentials(ctx, updatedData)
	if errFromCredentialsUpdate != nil {
		return nil, errFromCredentialsUpdate
	}
	_, err := u.userRepository.UpdateUserCredentials(updatedData)
	if err != nil {
		return nil, err
	}

	return &domains.BaseMessageResponse{
		Message: "You have successfully updated your credentials, please log in again.",
	}, nil
}

func (u UserService) DeleteUserById(id string) (*domains.UserDTO, *errors.ErrorStruct) {
	if id == "" {
		return nil, errors.NewError("Invalid ID format", 400)
	}
	user, err := u.userRepository.DeleteUserById(id)
	if err != nil {
		return nil, err
	}
	return &domains.UserDTO{
		Username:  user.Username,
		Email:     user.Email,
		ID:        id,
		Role:      user.Role,
		Confirmed: user.Confirmed,
	}, nil
}