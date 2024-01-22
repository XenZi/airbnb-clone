package services

import (
	"recommendation-service/domains"
	"recommendation-service/errors"
	"recommendation-service/repository"
	"time"
)

type RecommendationService struct {
	repo *repository.RecommendationRepository
}

func NewRecommendationService(repo *repository.RecommendationRepository) *RecommendationService {
	return &RecommendationService{
		repo: repo,
	}
}

func (rs RecommendationService) GetAllRecommendationsByUserID(id string) ([]domains.Recommendation, *errors.ErrorStruct) {
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	formattedDate := threeMonthsAgo.Format("2006-01-02")
	recommendations, err := rs.repo.GetAllRecommendationsForUser(id, formattedDate)
	if err != nil {
		return nil, err
	}
	return recommendations, nil
}
