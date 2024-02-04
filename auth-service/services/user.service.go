package services

import (
	"auth-service/client"
	"auth-service/config"
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/repository"
	"auth-service/utils"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	userRepository     *repository.UserRepository
	passwordService    *PasswordService
	jwtService         *JwtService
	validator          *utils.Validator
	encryptionService  *EncryptionService
	mailClient         client.MailClientInterface
	userClient         *client.UserClient
	notificationClient *client.NotificationClient
	tracer             trace.Tracer
	logger             *config.Logger
}

func NewUserService(userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JwtService,
	validator *utils.Validator,
	encryptionService *EncryptionService,
	mailClient client.MailClientInterface,
	userClient *client.UserClient,
	notificationClient *client.NotificationClient,
	tracer trace.Tracer,
	logger *config.Logger) *UserService {
	return &UserService{
		userRepository:     userRepo,
		passwordService:    passwordService,
		jwtService:         jwtService,
		validator:          validator,
		encryptionService:  encryptionService,
		mailClient:         mailClient,
		userClient:         userClient,
		notificationClient: notificationClient,
		tracer:             tracer,
		logger:             logger,
	}
}
func (u *UserService) CreateUser(ctx context.Context, registerUser domains.RegisterUser) (*domains.UserDTO, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()
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
		u.logger.LogError("user-service", fmt.Sprintf("Bad password for user %v", registerUser.Username))
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
		u.logger.LogError("user-service", fmt.Sprintf("Error occured while hashing the password"))
		return nil, errors.NewError(err.Error(), 500)
	}
	u.logger.LogInfo("user-service", "Password for user with username "+user.Username+" sucessfully hashed.")
	user.Password = hashedPassword
	newUser, foundErr := u.userRepository.SaveUser(ctx, user)
	if foundErr != nil {
		return nil, foundErr
	}
	id, err := newUser.ID.MarshalJSON()
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	token, encError := u.encryptionService.GenerateToken(string(id[1 : len(id)-1]))
	if encError != nil {
		u.logger.LogError("user-service", fmt.Sprintf("Error generating encrypted token for user with username %v", registerUser.Username))
		return nil, encError
	}
	u.logger.LogInfo("user-service", "Token for user with username "+user.Username+" sucessfully created.")

	errFromUserService := u.userClient.SendCreatedUser(ctx, newUser.ID.Hex(), registerUser)
	u.logger.LogInfo("user-service", "User with username "+user.Username+" sucessfully sent to user service for creation.")

	if errFromUserService != nil {
		_, err := u.userRepository.DeleteUserById(newUser.ID.Hex())
		if err != nil {
			u.logger.LogError("user-service", fmt.Sprintf("Error creating user with username %v at user-service, rolling back.", registerUser.Username))
			return nil, err
		}
		return nil, errFromUserService
	}
	u.logger.LogInfo("user-service", "User with username "+user.Username+" sucessfully created.")
	go func() {
		u.mailClient.SendAccountConfirmationEmail(registerUser.Email, token)
	}()
	errFromNotificationService := u.notificationClient.CreateNewUserStructNotification(ctx, newUser.ID.Hex())
	u.logger.LogInfo("user-service", "Creating user struct for user with username  "+user.Username+" in notification service.")

	if errFromNotificationService != nil {
		u.logger.LogError("user-service", fmt.Sprintf("Error creating user with username %v at notification-service.", registerUser.Username))
		u.logger.LogError("user-service", err.Error())
	}
	return &domains.UserDTO{
		ID:       string(id[1 : len(id)-1]),
		Email:    registerUser.Email,
		Username: registerUser.Username,
		Role:     registerUser.Role,
	}, nil
}

func (u *UserService) LoginUser(ctx context.Context, loginData domains.LoginUser) (*domains.SuccessfullyLoggedUser, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.LoginUser")
	defer span.End()
	u.logger.LogInfo("user-service", fmt.Sprintf("Looking for a user with email in our database %v", loginData.Email))

	user, err := u.userRepository.FindUserByEmail(ctx, loginData.Email)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User with email %v found", loginData.Email))

	isSamePassword := u.passwordService.CheckPasswordHash(loginData.Password, user.Password)
	if !isSamePassword {
		u.logger.LogError("user-service", "Bad credentials for user "+loginData.Email)
		return nil, errors.NewError("Bad credentials", 401)
	}
	jwtToken, foundError := u.jwtService.CreateKey(user.Email, user.Role, user.ID.Hex())
	if foundError != nil {
		u.logger.LogError("user-service", foundError.Error())
		return nil, errors.NewError(
			foundError.Error(),
			500)
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User with email %v successfully logged in the system.", loginData.Email))
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
	u.logger.LogInfo("user-service", fmt.Sprintf("User with ID %v tries to verify his account", userID))

	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	updatedUser, err := u.userRepository.UpdateUserConfirmation(userID)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User with ID %v verifies his account", userID))

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
	u.logger.LogInfo("user-service", fmt.Sprintf("User with email %v wants to reset a password", email))

	user, err := u.userRepository.FindUserByEmail(context.Background(), email)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, errors.NewError(err.GetErrorMessage(), err.GetErrorStatus())
	}
	token, err := u.encryptionService.GenerateToken(user.ID.Hex())
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, errors.NewError(err.GetErrorMessage(), err.GetErrorStatus())
	}
	go func() {
		u.mailClient.SendRequestResetPassword(email, token)
	}()
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v will get email", email))
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
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v wants to confirm password", userID))
	if requestData.Password != requestData.ConfirmedPassword {
		u.logger.LogError("user-service", "Password doesn't match")
		return nil, errors.NewError("Password doesn't match", 400)
	}
	if u.passwordService.CheckPasswordExistanceInBlacklist(requestData.Password) != false {
		u.logger.LogError("user-service", "Choose better password that is more secure!")
		return nil, errors.NewError("Choose better password that is more secure!", 400)
	}
	hashedPassword, hashError := u.passwordService.HashPassword(requestData.Password)
	if hashError != nil {
		return nil, err
	}

	user, err := u.userRepository.UpdateUserPassword(userID, hashedPassword)
	if err != nil {
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v confirmed new password", userID))
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
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v wants to change password", userID))
	user, err := u.userRepository.FindUserById(userID)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	if u.passwordService.CheckPasswordHash(data.OldPassword, user.Password) == false {
		u.logger.LogError("user-service", fmt.Sprintf("User %v wants to change password but it doesn't match with the old", userID))
		return nil, errors.NewError("Old password doesn't match", 400)
	}

	hashedPassword, hashError := u.passwordService.HashPassword(data.Password)
	if hashError != nil {
		return nil, errors.NewError(hashError.Error(), 500)
	}
	_, err = u.userRepository.UpdateUserPassword(userID, hashedPassword)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v changed the password.", userID))

	response := domains.BaseMessageResponse{
		Message: "You have updated your password",
	}
	return &response, nil
}

// OVDE SAM STAO ZA LOGOCE DA ZNAM
func (u UserService) UpdateCredentials(ctx context.Context, id string, updatedData domains.User) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	if updatedData.Email == "" || updatedData.Username == "" {
		return nil, errors.NewError("Email or username are empty", 400)
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v wants to update his credentials.", updatedData.ID.Hex()))

	foundUser, _ := u.userRepository.FindUserById(id)
	if foundUser.Username == updatedData.Username && foundUser.Email == updatedData.Email {
		u.logger.LogInfo("user-service", fmt.Sprintf("User %v updated his credentials.", updatedData.ID.Hex()))

		return &domains.BaseMessageResponse{
			Message: "You have successfully updated your credentials",
		}, nil
	}
	if foundUser.Email != updatedData.Email {
		_, err := u.userRepository.FindUserByEmail(context.Background(), updatedData.Email)
		if err == nil {
			u.logger.LogError("user-service", err.GetErrorMessage())
			return nil, errors.NewError("User with same email already exists", 400)
		}
	}
	if foundUser.Username != updatedData.Username {
		_, err := u.userRepository.FindUserByUsername(updatedData.Username)
		if err == nil {
			u.logger.LogError("user-service", err.GetErrorMessage())
			return nil, errors.NewError("User with same username already exists", 400)
		}
	}
	if u.passwordService.CheckPasswordHash(updatedData.Password, foundUser.Password) == false {
		u.logger.LogError("user-service", "Password doesn't match")
		return nil, errors.NewError("Password doesn't match", 400)
	}
	objectID, newError := primitive.ObjectIDFromHex(id)
	if newError != nil {
		u.logger.LogError("user-service", newError.Error())
		return nil, errors.NewError(newError.Error(), 500)
	}
	updatedData.ID = objectID
	errFromCredentialsUpdate := u.userClient.SendUpdateCredentials(ctx, updatedData)
	if errFromCredentialsUpdate != nil {
		u.logger.LogError("user-service", errFromCredentialsUpdate.GetErrorMessage())
		return nil, errFromCredentialsUpdate
	}
	_, err := u.userRepository.UpdateUserCredentials(updatedData)
	if err != nil {
		u.logger.LogError("user-service", err.GetErrorMessage())
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v updated his credentials.", updatedData.ID.Hex()))
	return &domains.BaseMessageResponse{
		Message: "You have successfully updated your credentials, please log in again.",
	}, nil
}

func (u UserService) DeleteUserById(id string) (*domains.UserDTO, *errors.ErrorStruct) {
	if id == "" {
		return nil, errors.NewError("Invalid ID format", 400)
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v wants to delete his account.", id))

	user, err := u.userRepository.DeleteUserById(id)
	if err != nil {
		return nil, err
	}
	u.logger.LogInfo("user-service", fmt.Sprintf("User %v deleted his account.", id))

	return &domains.UserDTO{
		Username:  user.Username,
		Email:     user.Email,
		ID:        id,
		Role:      user.Role,
		Confirmed: user.Confirmed,
	}, nil
}
