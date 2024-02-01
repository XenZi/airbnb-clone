package services

import (
	"context"
	"recommendation-service/client"
	"recommendation-service/domains"
	"recommendation-service/errors"
	"recommendation-service/repository"
	"time"
)

type RecommendationService struct {
	repo                *repository.RecommendationRepository
	accommodationClient *client.AccommodationClient
}

func NewRecommendationService(repo *repository.RecommendationRepository, accommodationClient *client.AccommodationClient) *RecommendationService {
	return &RecommendationService{
		repo:                repo,
		accommodationClient: accommodationClient,
	}
}

func (rs RecommendationService) GetAllRecommendationsByUserID(ctx context.Context, id string) ([]domains.AccommodationDTO, *errors.ErrorStruct) {
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	formattedDate := threeMonthsAgo.Format("2006-01-02")
	recommendations, err := rs.repo.GetAllRecommendationsForUser(id, formattedDate)
	if len(recommendations) == 0 {
		recommendations, err = rs.repo.GetAllRecommendationsByRating()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	returnedData, err := rs.accommodationClient.GetAllRecommendedAccommodationData(ctx, recommendations)
	if err != nil {
		return nil, err
	}

	return returnedData, nil
}

func (rs RecommendationService) GetAllRecommendationsByRating(ctx context.Context) ([]domains.AccommodationDTO, *errors.ErrorStruct) {
	recommendations, err := rs.repo.GetAllRecommendationsByRating()
	if err != nil {
		return nil, err
	}
	returnedData, err := rs.accommodationClient.GetAllRecommendedAccommodationData(ctx, recommendations)
	if err != nil {
		return nil, err
	}
	return returnedData, nil
}
