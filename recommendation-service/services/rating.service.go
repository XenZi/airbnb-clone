package services

import (
	"context"
	"recommendation-service/client"
	"recommendation-service/domains"
	"recommendation-service/errors"
	"recommendation-service/repository"
	"strings"
	"time"
)

type RatingService struct {
	repo                *repository.RatingRepository
	accommodationClient *client.AccommodationClient
	userClient          *client.UserClient
}

func NewRatingService(repo *repository.RatingRepository, accommodationClient *client.AccommodationClient, userClient *client.UserClient) *RatingService {
	return &RatingService{
		repo:                repo,
		accommodationClient: accommodationClient,
		userClient:          userClient,
	}
}

func (rs RatingService) CreateRatingForAccommodation(ctx context.Context, rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	rating.CreatedAt = strings.Split(time.Now().Local().String(), " ")[0]
	resp, err := rs.repo.RateAccommodation(rating)
	if err != nil {
		return nil, err
	}
	err = rs.accommodationClient.SendNewRatingForAccommodation(ctx, resp.AvgRating, resp.AccommodationID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) UpdateRatingForAccommodation(ctx context.Context, rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	rating.CreatedAt = strings.Split(time.Now().Local().String(), " ")[0]
	resp, err := rs.repo.UpdateRatingByAccommodationGuest(rating)
	if err != nil {
		return nil, err
	}
	err = rs.accommodationClient.SendNewRatingForAccommodation(ctx, resp.AvgRating, resp.AccommodationID)
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

func (rs RatingService) DeleteRatingForAccommodation(ctx context.Context, accommodationID, guestID string) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	data, err := rs.repo.DeleteRatingByGuestAndAccommodation(accommodationID, guestID)
	if err != nil {
		return nil, err
	}
	err = rs.accommodationClient.SendNewRatingForAccommodation(ctx, data.AvgRating, data.AccommodationID)

	return &domains.BaseMessageResponse{
		Message: "You have successfully deleted your rating for accommodation",
	}, nil
}

func (rs RatingService) CreateRatingForHost(ctx context.Context, rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	rating.CreatedAt = strings.Split(time.Now().Local().String(), " ")[0]
	resp, err := rs.repo.RateHost(rating)
	if err != nil {
		return nil, err
	}
	err = rs.userClient.SendNewRatingForUser(ctx, resp.AvgRating, rating.Host.ID)
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

func (rs RatingService) UpdateRatingForHostAndGuest(ctx context.Context, rateHost domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	rateHost.CreatedAt = strings.Split(time.Now().Local().String(), " ")[0]
	resp, err := rs.repo.UpdateRatingByHostAndGuest(rateHost)
	if err != nil {
		return nil, err
	}
	err = rs.userClient.SendNewRatingForUser(ctx, resp.AvgRating, rateHost.Host.ID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rs RatingService) DeleteRatingBetweenGuestAndHost(ctx context.Context, rateHost domains.RateHost) (*domains.BaseMessageResponse, *errors.ErrorStruct) {
	newAvgRating, err := rs.repo.DeleteRatingByHostAndUser(rateHost)
	if err != nil {
		return nil, err
	}
	err = rs.userClient.SendNewRatingForUser(ctx, newAvgRating, rateHost.Host.ID)
	if err != nil {
		return nil, err
	}
	return &domains.BaseMessageResponse{
		Message: "You have deleted your message successfully",
	}, nil
}

func (rs RatingService) GetRatingByGuestForAccommodation(guestID, accommodation string) (*domains.RateAccommodation, *errors.ErrorStruct) {
	rate, err := rs.repo.GetRatingByGuestForAccommodation(guestID, accommodation)
	if err != nil {
		return nil, err
	}
	return rate, nil
}
