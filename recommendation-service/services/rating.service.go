package services

import (
	"recommendation-service/domains"
	"recommendation-service/errors"
	"recommendation-service/repository"
	"time"
)

type RatingService struct {
	repo *repository.RatingRepository
}

func NewRatingService(repo *repository.RatingRepository) *RatingService {
	return &RatingService{
		repo: repo,
	}
}

func (rs RatingService) CreateRatingForAccommodation(rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	rating.CreatedAt = time.Now().Local().String()
	resp, err := rs.repo.RateAccommodation(rating)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) UpdateRatingForAccommodation(rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	rating.CreatedAt = time.Now().Local().String()
	resp, err := rs.repo.UpdateRatingByAccommodationGuest(rating)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) GetAllAccommodationRatings(id string) (*[]domains.RateAccommodation, *errors.ErrorStruct) {
	resp, err := rs.repo.GetAllRatingsByAccommodation(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) DeleteRatingForAccommodation(rating domains.RateAccommodation) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	err := rs.repo.DeleteRatingByGuestAndAccommodation(rating)
	if err != nil {
		return nil, err
	}
	return &domains.BaseMessageResponse{
		Message: "You have successfully deleted your rating for accommodation",
	}, nil
}

func (rs RatingService) CreateRatingForHost(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	rating.CreatedAt = time.Now().Local().String()
	resp, err := rs.repo.RateHost(rating)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) GetAllRatingsForHostByID(id string) (*[]domains.RateHost, *errors.ErrorStruct) {
	resp, err := rs.repo.GetAllRatingsByHostID(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) UpdateRatingForHostAndGuest(rateHost domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	rateHost.CreatedAt = time.Now().Local().String()
	resp, err := rs.repo.UpdateRatingByHostAndGuest(rateHost)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) DeleteRatingBetweenGuestAndHost(rateHost domains.RateHost) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	err := rs.repo.DeleteRatingByHostAndUser(rateHost)
	if err != nil {
		return nil, err
	}
	return &domains.BaseMessageResponse{
		Message: "You have deleted your message successfully",
	}, nil
}
