package repository

import (
	"context"
	"log"
	"recommendation-service/domains"
	"recommendation-service/errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type RecommendationRepository struct {
	driver neo4j.DriverWithContext
}

func NewRecommendationRepository(driver neo4j.DriverWithContext) *RecommendationRepository {
	return &RecommendationRepository{
		driver: driver,
	}
}

func (rr RecommendationRepository) GetAllRecommendationsForUser(id string, pastThreeMontsDate string) ([]domains.Recommendation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := rr.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)
	recommendedAccommodations, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (u:Guest {id: $guestID})-[r1:RATED]->(a:Accommodation)
				WITH u, collect(a) AS accommodationsRatedByUser
				MATCH (similarUser:Guest)-[r2:RATED]->(accommodation:Accommodation)
				WHERE similarUser <> u AND r2.rate > 3.5
				WITH u, accommodationsRatedByUser, collect(DISTINCT accommodation) AS accommodationsAbove
				UNWIND accommodationsAbove AS a
				MATCH (a)<-[ratings:RATED]-(guest:Guest)
				WHERE ratings.rate >= 1 AND date(ratings.createdAt) >= date($passedDate)
				WITH a, count(ratings) AS lowRatingsCount
				WHERE lowRatingsCount < 5
				WITH a
				MATCH (a)<-[r:RATED]-(guest:Guest)
				WITH a, avg(r.rate) AS avgRating
				RETURN a.id AS accommodationId, avgRating
				ORDER BY avgRating DESC`,
				map[string]any{
					"guestID":    id,
					"passedDate": pastThreeMontsDate,
				})
			if err != nil {
				return nil, err
			}
			var recommendations []domains.Recommendation
			for result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationId")
				rating, _ := record.Get("avgRating")
				recc := domains.Recommendation{
					AccommodationID: accommodationID.(string),
					Rating:          rating.(float64),
				}
				recommendations = append(recommendations, recc)
			}

			return recommendations, nil
		})

	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return recommendedAccommodations.([]domains.Recommendation), nil
}

func (rr RecommendationRepository) GetAllRecommendationsByRating() ([]domains.Recommendation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := rr.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)
	recommendedAccommodations, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (a:Accommodation)<-[r:RATED]-()
				WITH a, AVG(r.rate) AS averageRating
				WHERE averageRating > 3.5
				RETURN a.id AS accommodationId, averageRating AS avgRating
				ORDER BY avgRating DESC
				LIMIT 10				
				`,
				map[string]any{})
			if err != nil {
				return nil, err
			}
			var recommendations []domains.Recommendation
			for result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationId")
				log.Println(accommodationID)
				rating, _ := record.Get("avgRating")
				recc := domains.Recommendation{
					AccommodationID: accommodationID.(string),
					Rating:          rating.(float64),
				}
				recommendations = append(recommendations, recc)
			}

			return recommendations, nil
		})

	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	return recommendedAccommodations.([]domains.Recommendation), nil
}
