package repository

import (
	"context"
	"log"
	"recommendation-service/domains"
	"recommendation-service/errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type RatingRepository struct {
	driver neo4j.DriverWithContext
}

func NewRatingRepository(driver neo4j.DriverWithContext) *RatingRepository {
	return &RatingRepository{
		driver: driver,
	}
}

func (r RatingRepository) RateAccommodation(rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	rateAccommodation, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MERGE (g:Guest {id: $guestID})
				ON CREATE SET g.email = $guestEmail, g.username = $guestUsername
			MERGE (a:Accommodation {id: $accommodationID})
				ON CREATE SET a.id = $accommodationID
			MERGE (g)-[r:RATED]->(a)
				ON CREATE SET r.rate = $rate, r.createdAt = $createdAt
			WITH a, g, r
			MATCH (a)<-[r2:RATED]-()
			WITH a, g, r, AVG(r2.rate) AS averageRating
			SET a.averageRating = averageRating
			
			RETURN a.id as accommodationId, a.averageRating AS avgRating,
				   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
				   r.rate AS rate, r.createdAt AS createdAt`,
				map[string]any{
					"guestID":         rating.Guest.ID,
					"guestEmail":      rating.Guest.Email,
					"guestUsername":   rating.Guest.Username,
					"accommodationID": rating.AccommodationID,
					"rate":            rating.Rate,
					"createdAt":       rating.CreatedAt,
				})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationId")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				avgRating, _ := record.Get("avgRating")
				rate, _ := record.Get("rate")
				createdAt, _ := record.Get("createdAt")
				rateAccommodation := domains.RateAccommodation{
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					AccommodationID: accommodationID.(string),
					Rate:            rate.(int64),
					CreatedAt:       createdAt.(string),
					AvgRating:       avgRating.(float64),
				}
				return rateAccommodation, nil
			}

			return nil, result.Err()
		})

	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	rateAccResult := rateAccommodation.(domains.RateAccommodation)
	return &rateAccResult, nil
}

func (r RatingRepository) UpdateRatingByAccommodationGuest(rating domains.RateAccommodation) (*domains.RateAccommodation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)
	log.Println(rating)
	rateAccommodation, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (a:Accommodation {id: $accommodationId})<-[r:RATED]-(g:Guest {id: $guestId})
				SET r.rate = $newRate
				
				WITH a, g, r
				MATCH (a)<-[r2:RATED]-()
				WITH a, g, r, AVG(r2.rate) AS averageRating
				SET a.averageRating = averageRating
				
				RETURN a.id AS accommodationID, a.averageRating AS avgRating,
					   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
					   r.rate AS rate, r.createdAt AS createdAt`,
				map[string]any{
					"guestId":         rating.Guest.ID,
					"accommodationId": rating.AccommodationID,
					"newRate":         rating.Rate,
					"createdAt":       rating.CreatedAt,
				})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationID")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				createdAt, _ := record.Get("createdAt")
				rate, _ := record.Get("rate")
				avgRating, _ := record.Get("avgRating")
				rateAccommodation := domains.RateAccommodation{
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					AccommodationID: accommodationID.(string),
					Rate:            rate.(int64),
					CreatedAt:       createdAt.(string),
					AvgRating:       avgRating.(float64),
				}
				return rateAccommodation, nil
			}
			return nil, result.Err()
		})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	rateAccommodationResult := rateAccommodation.(domains.RateAccommodation)
	return &rateAccommodationResult, nil
}

func (r RatingRepository) DeleteRatingByGuestAndAccommodation(accommodationID, guestID string) (*domains.RateAccommodation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	data, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (a:Accommodation {id: $accommodationId})<-[r:RATED]-(g:Guest {id: $guestId})
				DETACH DELETE r
				WITH a
				OPTIONAL MATCH (a)<-[r2:RATED]-()
				WITH a, COALESCE(AVG(r2.rate), 0.0) AS averageRating
				SET a.averageRating = averageRating
				RETURN a.id AS accommodationID, a.averageRating as avgRating`,
				map[string]any{
					"guestId":         guestID,
					"accommodationId": accommodationID,
				})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				record := result.Record()
				avgRating, _ := record.Get("avgRating")
				accommodationID, _ := record.Get("accommodationID")
				rateAccommodation := domains.RateAccommodation{
					AccommodationID: accommodationID.(string),
					AvgRating:       avgRating.(float64),
				}
				return rateAccommodation, nil
			}
			return nil, result.Err()
		})

	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}

	if data == nil {
		return nil, errors.NewError("Resource not found", 404)
	}
	rateAccommodationResult := data.(domains.RateAccommodation)
	return &rateAccommodationResult, nil
}

func (r RatingRepository) GetAllRatingsByHostID(hostID string) (*[]domains.RateHost, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	ratingResults, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (h:Host {id: $hostID})<-[r:RATED]-(g:Guest)
				WITH h, g, r
				MATCH (otherGuest:Guest)-[otherRating:RATED]->(h)
				WITH h, g, r, AVG(otherRating.rate) AS avgRating
				RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, 
					   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
					   r.rate AS rate, r.createdAt AS createdAt, avgRating
				
			`,
				map[string]any{
					"hostID": hostID,
				})
			if err != nil {
				return nil, err
			}
			var hostRatings []domains.RateHost
			for result.Next(ctx) {
				record := result.Record()
				hostID, _ := record.Get("hostID")
				hostEmail, _ := record.Get("hostEmail")
				hostUsername, _ := record.Get("hostUsername")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				createdAt, _ := record.Get("createdAt")
				rate, _ := record.Get("rate")
				avgRating, _ := record.Get("avgRating")
				rateHost := domains.RateHost{
					Host: domains.Host{
						ID:       hostID.(string),
						Email:    hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate:      rate.(int64),
					CreatedAt: createdAt.(string),
					AvgRating: avgRating.(float64),
				}
				hostRatings = append(hostRatings, rateHost)
			}
			if result.Err() != nil {
				return nil, result.Err()
			}
			return hostRatings, nil
		})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	rateHostResult := ratingResults.([]domains.RateHost)
	return &rateHostResult, nil
}

func (r RatingRepository) RateHost(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	rateHost, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MERGE (g:Guest {id: $guestID})
				ON CREATE SET g.email = $guestEmail, g.username = $guestUsername
			  MERGE (h:Host {id: $hostID})
				ON CREATE SET h.email = $hostEmail, h.username = $hostUsername
			  MERGE (g)-[r:RATED]->(h)
				ON CREATE SET r.rate = $rate, r.createdAt = $createdAt
			  WITH h, g, r.rate AS rate, r.createdAt AS createdAt
			  MATCH (otherGuest:Guest)-[otherRating:RATED]->(h)
			  WITH h, g, rate, createdAt, AVG(otherRating.rate) AS avgRating
			  SET h.avgRating = avgRating
			  RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername,
					 g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername,
					 createdAt, rate, avgRating`,
				map[string]any{
					"guestID":       rating.Guest.ID,
					"guestEmail":    rating.Guest.Email,
					"guestUsername": rating.Guest.Username,
					"hostID":        rating.Host.ID,
					"hostEmail":     rating.Host.Email,
					"hostUsername":  rating.Host.Username,
					"rate":          rating.Rate,
					"createdAt":     rating.CreatedAt,
				})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				record := result.Record()
				hostID, _ := record.Get("hostID")
				hostEmail, _ := record.Get("hostEmail")
				hostUsername, _ := record.Get("hostUsername")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				createdAt, _ := record.Get("createdAt")
				rate, _ := record.Get("rate")
				avgRating, _ := record.Get("avgRating")
				rateHost := domains.RateHost{
					Host: domains.Host{
						ID:       hostID.(string),
						Email:    hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate:      rate.(int64),
					CreatedAt: createdAt.(string),
					AvgRating: avgRating.(float64),
				}

				return rateHost, nil
			}

			return nil, result.Err()
		})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	rateHostResult := rateHost.(domains.RateHost)
	return &rateHostResult, nil
}

func (r RatingRepository) GetAllRatingsByAccommodation(accommodationID string) (*[]domains.RateAccommodation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	ratingResults, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (a:Accommodation {id: $accommodationID})<-[r:RATED]-(g:Guest)
				RETURN a.id as accommodationID, 
					g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
					r.rate AS rate, r.createdAt as createdAt
				`,
				map[string]any{
					"accommodationID": accommodationID,
				})
			if err != nil {
				return nil, err
			}
			var accommodationRatings []domains.RateAccommodation
			for result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationID")

				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")

				rate, _ := record.Get("rate")
				createdAt, _ := record.Get("createdAt")

				accommodationHost := domains.RateAccommodation{
					AccommodationID: accommodationID.(string),
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate:      rate.(int64),
					CreatedAt: createdAt.(string),
				}
				accommodationRatings = append(accommodationRatings, accommodationHost)
			}
			if result.Err() != nil {
				return nil, result.Err()
			}
			return accommodationRatings, nil
		})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	accommodationRatingResults := ratingResults.([]domains.RateAccommodation)
	return &accommodationRatingResults, nil
}

func (r RatingRepository) UpdateRatingByHostAndGuest(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	rateHost, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (h:Host {id: $hostId})<-[r:RATED]-(g:Guest {id: $guestId})
				SET r.rate = $newRate, r.createdAt = $createdAt
				WITH h
				MATCH (otherGuest:Guest)-[otherRating:RATED]->(h)
				WITH h, AVG(otherRating.rate) AS avgRating
				SET h.avgRating = avgRating
				RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, 
					   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
					   r.rate AS rate, h.avgRating AS avgRating
				`,
				map[string]any{
					"guestId":   rating.Guest.ID,
					"hostId":    rating.Host.ID,
					"newRate":   rating.Rate,
					"createdAt": rating.CreatedAt,
				})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				record := result.Record()
				hostID, _ := record.Get("hostID")
				hostEmail, _ := record.Get("hostEmail")
				hostUsername, _ := record.Get("hostUsername")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				avgRating, _ := record.Get("avgRating")
				rate, _ := record.Get("rate")
				createdAt, _ := record.Get("createdAt")

				rateHost := domains.RateHost{
					Host: domains.Host{
						ID:       hostID.(string),
						Email:    hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate:      rate.(int64),
					CreatedAt: createdAt.(string),
					AvgRating: avgRating.(float64),
				}
				return rateHost, nil
			}
			return nil, result.Err()
		})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}

	rateHostResult := rateHost.(domains.RateHost)
	return &rateHostResult, nil
}

func (r RatingRepository) DeleteRatingByHostAndUser(rating domains.RateHost) (float64, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	returnedVal, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (h:Host {id: $hostId})<-[r:RATED]-(g:Guest {id: $guestId})
				DETACH DELETE r
				WITH h
				MATCH (otherGuest:Guest)-[otherRating:RATED]->(h)
				WITH h, COALESCE(AVG(otherRating.rate), 0) AS avgRating
				SET h.avgRating = avgRating
				RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, h.avgRating AS avgRating`,
				map[string]any{
					"guestId": rating.Guest.ID,
					"hostId":  rating.Host.ID,
				})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				record := result.Record()
				avgRating, _ := record.Get("avgRating")
				return avgRating.(float64), nil
			}
			return nil, result.Err()
		})

	if err != nil {
		return -1, errors.NewError(err.Error(), 500)
	}
	return returnedVal.(float64), nil
}

func (r RatingRepository) GetRatingByGuestForAccommodation(guestID, accommodationID string) (*domains.RateAccommodation, *errors.ErrorStruct) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	rateAccommodation, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				`MATCH (g:Guest {id: $guestID})-[r:RATED]->(a:Accommodation {id: $accommodationID})
                WITH a, g, r
                MATCH (a)<-[r2:RATED]-()
                WITH a, g, r, AVG(r2.rate) AS averageRating
                RETURN a.id AS accommodationId, averageRating AS avgRating,
                       g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername,
                       r.rate AS rate, r.createdAt AS createdAt
				`,
				map[string]any{
					"guestID":         guestID,
					"accommodationID": accommodationID,
				})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				record := result.Record()
				accommodationID, _ := record.Get("accommodationId")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
				avgRating, _ := record.Get("avgRating")
				rate, _ := record.Get("rate")
				createdAt, _ := record.Get("createdAt")
				rateAccommodation := domains.RateAccommodation{
					Guest: domains.Guest{
						ID:       guestID.(string),
						Email:    guestEmail.(string),
						Username: guestUsername.(string),
					},
					AccommodationID: accommodationID.(string),
					Rate:            rate.(int64),
					CreatedAt:       createdAt.(string),
					AvgRating:       avgRating.(float64),
				}
				return rateAccommodation, nil
			}

			return nil, result.Err()
		})

	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}

	if rateAccommodation == nil {
		return nil, errors.NewError("Resource not found", 404)
	}
	rateAccResult := rateAccommodation.(domains.RateAccommodation)
	return &rateAccResult, nil
}
