package services

import (
	"recommendation-service/domains"
	"recommendation-service/errors"
	"recommendation-service/repository"
)

type RatingService struct {
	repo *repository.RatingRepository
}

func NewRatingService(repo *repository.RatingRepository) *RatingService {
	return &RatingService{
		repo: repo,
	}
}


func (rs RatingService) CreateRatingForAccommodation(rating domains.RateAccommodation) {
	rs.repo.RateAccommodation(rating)
}

func (rs RatingService) CreateRatingForHost(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
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
	resp, err := rs.repo.UpdateRatingByHostAndGuest(rateHost)
	if err != nil {
		return nil, err
	}
	return resp, nil
}