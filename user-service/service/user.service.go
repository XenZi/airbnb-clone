package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-service/config"

	"user-service/client"
	"user-service/domain"
	"user-service/errors"
	"user-service/repository"
	"user-service/utils"
)

type UserService struct {
	userRepository    *repository.UserRepository
	validator         *utils.Validator
	reservationClient *client.ReservationClient
	authClient        *client.AuthClient
	accClient         *client.AccClient
	logger            *config.Logger
}

const source = "user-service"

func NewUserService(userRepo *repository.UserRepository,
	validator *utils.Validator,
	reservationClient *client.ReservationClient,
	authClient *client.AuthClient,
	accClient *client.AccClient,
	logger *config.Logger) *UserService {
	return &UserService{
		userRepository:    userRepo,
		validator:         validator,
		reservationClient: reservationClient,
		authClient:        authClient,
		accClient:         accClient,
		logger:            logger,
	}
}

func (u *UserService) CreateUser(createUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	u.validator.ValidateUser(&createUser)
	validErrors := u.validator.GetErrors()
	if len(validErrors) > 0 {
		var constructedError string
		for _, message := range validErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	foundId, erro := primitive.ObjectIDFromHex(createUser.ID)
	if erro != nil {
		message := erro.Error()
		u.logger.LogError(source, message)
		return nil, errors.NewError(message, 500)
	}
	user := domain.User{
		ID:        foundId,
		Username:  createUser.Username,
		Email:     createUser.Email,
		Role:      createUser.Role,
		FirstName: createUser.FirstName,
		LastName:  createUser.LastName,
		Residence: createUser.Residence,
		Age:       createUser.Age,
		Rating:    createUser.Rating,
	}
	newUser, foundErr := u.userRepository.CreatUser(user)
	if foundErr != nil {
		u.logger.LogError(source, foundErr.GetErrorMessage())
		return nil, foundErr
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by ID: %v created", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) UpdateUser(updateUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	u.validator.ValidateUser(&updateUser)
	validErrors := u.validator.GetErrors()
	if len(validErrors) > 0 {
		var constructedError string
		for _, message := range validErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	foundId, err := primitive.ObjectIDFromHex(updateUser.ID)
	if err != nil {
		message := err.Error()
		u.logger.LogError(source, message)
		return nil, errors.NewError(message, 500)
	}
	user := domain.User{
		ID:        foundId,
		Username:  updateUser.Username,
		Email:     updateUser.Email,
		Role:      updateUser.Role,
		FirstName: updateUser.FirstName,
		LastName:  updateUser.LastName,
		Residence: updateUser.Residence,
		Age:       updateUser.Age,
		Rating:    updateUser.Rating,
	}
	newUser, foundErr := u.userRepository.UpdateUser(user)
	if foundErr != nil {
		u.logger.LogError(source, foundErr.GetErrorMessage())
		return nil, foundErr
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by ID: %v updated", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) UpdateUserCreds(updateUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	u.validator.ValidateCreds(&updateUser)
	validErrors := u.validator.GetErrors()
	if len(validErrors) > 0 {
		var constructedError string
		for _, message := range validErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	foundId, err := primitive.ObjectIDFromHex(updateUser.ID)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	user := domain.User{
		ID:       foundId,
		Username: updateUser.Username,
		Email:    updateUser.Email,
	}
	newUser, foundErr := u.userRepository.UpdateUserCreds(user)
	if foundErr != nil {
		return nil, foundErr
	}

	u.logger.LogInfo(source, fmt.Sprintf("User creds by ID: %v updated", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) GetAllUsers() ([]*domain.User, *errors.ErrorStruct) {
	userCollection, err := u.userRepository.GetAllUsers()
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return nil, err
	}
	return userCollection, nil
}

func (u *UserService) GetUserById(id string) (*domain.User, *domain.HostUser, *errors.ErrorStruct) {
	foundUser, err := u.userRepository.GetUserById(id)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return nil, nil, err
	}
	if foundUser.Role == "Host" {
		dist, erro := u.isDistinguished(foundUser.ID.Hex(), foundUser.Rating)
		if erro != nil {
			message := erro.GetErrorMessage()
			u.logger.LogError(source, message)
			dist = false
		}
		hostUser := domain.HostUser{
			ID:            foundUser.ID,
			Username:      foundUser.Username,
			Email:         foundUser.Email,
			Role:          foundUser.Role,
			FirstName:     foundUser.FirstName,
			LastName:      foundUser.LastName,
			Residence:     foundUser.Residence,
			Age:           foundUser.Age,
			Rating:        foundUser.Rating,
			Distinguished: dist,
		}
		u.logger.LogInfo(source, fmt.Sprintf("Host user by ID: %v found", hostUser.ID.Hex()))
		return nil, &hostUser, nil
	}
	u.logger.LogInfo(source, fmt.Sprintf("Guest user by ID: %v found", foundUser.ID.Hex()))
	return foundUser, nil, nil
}

func (u *UserService) isDistinguished(id string, rating float64) (bool, *errors.ErrorStruct) {
	reqs, err := u.reservationClient.CheckDistinguished(context.TODO(), id)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return false, err
	}
	if reqs && rating > 4.0 {
		u.logger.LogInfo(source, fmt.Sprintf("Host user by ID: %v is distinguished", id))
		return true, nil
	}
	u.logger.LogInfo(source, fmt.Sprintf("Host user by ID: %v is NOT distinguished", id))
	return false, nil
}

func (u *UserService) DeleteUser(role string, id string) *errors.ErrorStruct {
	err := u.reservationClient.UserDeleteAllowed(context.TODO(), id, role)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return err
	}
	if role == "Host" {
		err := u.accClient.DeleteUserAccommodations(context.TODO(), id)
		if err != nil {
			message := err.GetErrorMessage()
			u.logger.LogError(source, message)
			return err
		}
	}
	if role != "Guest" && role != "Host" {
		err2 := errors.NewError(fmt.Sprintf("User by id: %v not allowed by role: %s", id, role), 401)
		u.logger.LogError(source, err2.GetErrorMessage())
		return err2
	}
	newErr := u.authClient.DeleteUserAuth(context.TODO(), id)
	if newErr != nil {
		message := newErr.GetErrorMessage()
		u.logger.LogError(source, message)
		return newErr
	}
	err3 := u.userRepository.DeleteUser(id)
	if err3 != nil {
		message := err3.GetErrorMessage()
		u.logger.LogError(source, message)
		return err3
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by id: %v deleted", id))
	return nil
}

func (u *UserService) UpdateRating(id string, rating float64) *errors.ErrorStruct {
	err := u.userRepository.UpdateRating(id, rating)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return err
	}
	u.logger.LogInfo(source, fmt.Sprintf("User rating by id: %v updated", id))
	return nil
}
