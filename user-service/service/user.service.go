package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"log"
	"user-service/client"
	"user-service/domain"
	"user-service/errors"
	"user-service/repository"
	"user-service/utils"
)

type UserService struct {
	userRepository    *repository.UserRepository
	jwtService        *JwtService
	validator         *utils.Validator
	reservationClient *client.ReservationClient
	authClient        *client.AuthClient
	accClient         *client.AccClient
}

func NewUserService(userRepo *repository.UserRepository,
	jwtService *JwtService,
	validator *utils.Validator,
	reservationClient *client.ReservationClient,
	authClient *client.AuthClient,
	accClient *client.AccClient) *UserService {
	return &UserService{
		userRepository:    userRepo,
		jwtService:        jwtService,
		validator:         validator,
		reservationClient: reservationClient,
		authClient:        authClient,
		accClient:         accClient,
	}
}

func (u *UserService) CreateUser(createUser domain.CreateUser) (*domain.User, *errors.ErrorStruct) {
	log.Println(createUser)
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
		return nil, errors.NewError(erro.Error(), 500)
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
		return nil, foundErr
	}

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
		return nil, errors.NewError(err.Error(), 500)
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
		return nil, foundErr
	}

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

	return newUser, nil
}

func (u *UserService) GetAllUsers() ([]*domain.User, *errors.ErrorStruct) {
	userCollection, err := u.userRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return userCollection, nil
}

func (u *UserService) GetUserById(id string) (*domain.User, *domain.HostUser, *errors.ErrorStruct) {
	foundUser, err := u.userRepository.GetUserById(id)
	if err != nil {
		return nil, nil, err
	}
	if foundUser.Role == "Host" {
		dist := false
		dist, err = u.isDistinguished(foundUser.ID.Hex(), foundUser.Rating)
		if err != nil {
			return nil, nil, err
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
		return nil, &hostUser, nil
	}
	return foundUser, nil, nil
}

func (u *UserService) isDistinguished(id string, rating float64) (bool, *errors.ErrorStruct) {
	reqs, err := u.reservationClient.CheckDistinguished(context.TODO(), id)
	if err != nil {
		return false, err
	}
	// Req 4
	if reqs && rating > 4.0 {
		return true, nil
	}
	return false, nil
}

func (u *UserService) DeleteUser(role string, id string) *errors.ErrorStruct {
	allow, err := u.reservationClient.UserDeleteAllowed(context.TODO(), id, role)
	if err != nil {
		log.Println("ovo je error", err)
		return err
	}
	if !allow {
		return errors.NewError("user has reservations", 401)
	}
	if role == "Host" {
		err := u.accClient.DeleteUserAccommodations(context.TODO(), id)
		if err != nil {
			return err
		}
	}
	if role != "Guest" && role != "Host" {
		return errors.NewError("not allowed by role", 401)
	}
	newErr := u.authClient.DeleteUserAuth(context.TODO(), id)
	if newErr != nil {
		return errors.NewError("auth deletion error", 500)
	}
	erro := u.userRepository.DeleteUser(id)
	if erro != nil {
		return errors.NewError("internal error", 500)
	}
	return nil
}

func (u *UserService) UpdateRating(id string, rating float64) *errors.ErrorStruct {
	err := u.userRepository.UpdateRating(id, rating)
	if err != nil {
		return err
	}
	return nil
}
