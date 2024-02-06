package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
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
	tracer            trace.Tracer
}

const source = "user-service"

func NewUserService(userRepo *repository.UserRepository,
	validator *utils.Validator,
	reservationClient *client.ReservationClient,
	authClient *client.AuthClient,
	accClient *client.AccClient,
	logger *config.Logger,
	tracer trace.Tracer) *UserService {
	return &UserService{
		userRepository:    userRepo,
		validator:         validator,
		reservationClient: reservationClient,
		authClient:        authClient,
		accClient:         accClient,
		logger:            logger,
		tracer:            tracer,
	}
}

func (u *UserService) CreateUser(ctx context.Context, createUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()
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
	newUser, foundErr := u.userRepository.CreatUser(ctx, user)
	if foundErr != nil {
		u.logger.LogError(source, foundErr.GetErrorMessage())
		return nil, foundErr
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by ID: %v created", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) UpdateUser(ctx context.Context, updateUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()
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
	newUser, foundErr := u.userRepository.UpdateUser(ctx, user)
	if foundErr != nil {
		u.logger.LogError(source, foundErr.GetErrorMessage())
		return nil, foundErr
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by ID: %v updated", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) UpdateUserCreds(ctx context.Context, updateUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateUserCreds")
	defer span.End()
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
	newUser, foundErr := u.userRepository.UpdateUserCreds(ctx, user)
	if foundErr != nil {
		return nil, foundErr
	}

	u.logger.LogInfo(source, fmt.Sprintf("User creds by ID: %v updated", newUser.ID.Hex()))
	return newUser, nil
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.GetAllUsers")
	defer span.End()
	userCollection, err := u.userRepository.GetAllUsers(ctx)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return nil, err
	}
	return userCollection, nil
}

func (u *UserService) GetUserById(ctx context.Context, id string) (*domain.User, *domain.HostUser, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.GetUserById")
	defer span.End()
	foundUser, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return nil, nil, err
	}
	if foundUser.Role == "Host" {
		dist, erro := u.isDistinguished(ctx, foundUser.ID.Hex(), foundUser.Rating)
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

func (u *UserService) isDistinguished(ctx context.Context, id string, rating float64) (bool, *errors.ErrorStruct) {
	ctx, span := u.tracer.Start(ctx, "UserService.isDistinguished")
	defer span.End()
	reqs, err := u.reservationClient.CheckDistinguished(ctx, id)
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

func (u *UserService) DeleteUser(ctx context.Context, role string, id string) *errors.ErrorStruct {
	ctx, span := u.tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()
	err := u.reservationClient.UserDeleteAllowed(ctx, id, role)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return err
	}
	if role == "Host" {
		err := u.accClient.DeleteUserAccommodations(ctx, id)
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
	newErr := u.authClient.DeleteUserAuth(ctx, id)
	if newErr != nil {
		message := newErr.GetErrorMessage()
		u.logger.LogError(source, message)
		return newErr
	}
	err3 := u.userRepository.DeleteUser(ctx, id)
	if err3 != nil {
		message := err3.GetErrorMessage()
		u.logger.LogError(source, message)
		return err3
	}
	u.logger.LogInfo(source, fmt.Sprintf("User by id: %v deleted", id))
	return nil
}

func (u *UserService) UpdateRating(ctx context.Context, id string, rating float64) *errors.ErrorStruct {
	ctx, span := u.tracer.Start(ctx, "UserService.UpdateRating")
	defer span.End()
	err := u.userRepository.UpdateRating(ctx, id, rating)
	if err != nil {
		message := err.GetErrorMessage()
		u.logger.LogError(source, message)
		return err
	}
	u.logger.LogInfo(source, fmt.Sprintf("User rating by id: %v updated", id))
	return nil
}
