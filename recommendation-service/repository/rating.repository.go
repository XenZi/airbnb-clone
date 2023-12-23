package repository

import (
	"context"
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
func(r RatingRepository) CreateGuest(guest domains.Guest) error {
	ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
    defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, 
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (g:Guest) SET g.id = $id, g.email = $email, g.username = $username RETURN g.name + ', from node ' + id(g)",
				map[string]any{"id": guest.ID, "email": guest.Email, "username": guest.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})


    return err
}



func(r RatingRepository) CreateHost(host domains.Host) error {
	ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
    defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, 
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (h:Host) SET h.id = $id, h.email = $email, h.username = $username RETURN h.name + ', from node ' + id(g)",
				map[string]any{"id": host.ID, "email": host.Email, "username": host.Username})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})


    return err
}

func(r RatingRepository) CreateAccommodation(accommodation domains.Accommodation) error {
	ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
    defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, 
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"CREATE (a:Accommodation) SET a.id = $id RETURN id(a)",
				map[string]any{"id": accommodation.ID})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})


    return err
}

func(r RatingRepository) RateAccommodation(rating domains.RateAccommodation) error {
	ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
    defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, 
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
			`MERGE (g:Guest {id: $guestID})
			ON CREATE SET g.email = $guestEmail, g.username = $guestUsername
			MERGE (a:Accommodation {id: $accommodationID})
			ON CREATE SET a.id = $accommodationID
			MERGE (g)-[r:RATED]->(a)
			ON CREATE SET r.rate = $rate
			RETURN r`,
            map[string]any{                
			"guestID":           rating.Guest.ID,
			"guestEmail":        rating.Guest.Email,
			"guestUsername":     rating.Guest.Username,
			"accommodationID":   rating.AccommodationID,
			"rate":              rating.Rate,
			})			
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})

    return err
}

func(r RatingRepository) RateHost(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
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
			ON CREATE SET r.rate = $rate
			RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, 
			g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
			r.rate AS rate`,
            map[string]any{                
			"guestID":           rating.Guest.ID,
			"guestEmail":        rating.Guest.Email,
			"guestUsername":     rating.Guest.Username,
			"hostID":           rating.Host.ID,
			"hostEmail":        rating.Host.Email,
			"hostUsername":     rating.Host.Username,
			"rate":              rating.Rate,
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
			
				rate, _ := record.Get("rate")
			
				rateHost := domains.RateHost{
					Host: domains.Host{
						ID: hostID.(string),
						Email: hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID: guestID.(string),
						Email: guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate: rate.(int64),
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
			RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, 
				   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
				   r.rate AS rate
			`,
            map[string]any{                
				"hostID": hostID,
			})			
			if err != nil {
				return nil, err
			}
			var hostRatings []domains.RateHost
			if result.Next(ctx) {
				record := result.Record()
				hostID, _ := record.Get("hostID")
				hostEmail, _ := record.Get("hostEmail")
				hostUsername, _ := record.Get("hostUsername")
				guestID, _ := record.Get("guestID")
				guestEmail, _ := record.Get("guestEmail")
				guestUsername, _ := record.Get("guestUsername")
			
				rate, _ := record.Get("rate")
			
				rateHost := domains.RateHost{
					Host: domains.Host{
						ID: hostID.(string),
						Email: hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID: guestID.(string),
						Email: guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate: rate.(int64),
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

func (r RatingRepository) UpdateRatingByHostAndGuest(rating domains.RateHost) (*domains.RateHost, *errors.ErrorStruct) {
	ctx := context.Background()
    session := r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
    defer session.Close(ctx)

	rateHost, err := session.ExecuteWrite(ctx, 
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
			`MATCH (h:Host {id: $hostId})<-[r:RATED]-(g:Guest {id: $guestId}) SET r.rate = $newRate
			RETURN h.id AS hostID, h.email AS hostEmail, h.username AS hostUsername, 
				   g.id AS guestID, g.email AS guestEmail, g.username AS guestUsername, 
				   r.rate AS rate`,
            map[string]any{                
			"guestId":           rating.Guest.ID,
			"hostId":           rating.Host.ID,
			"newRate":          	rating.Rate,
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
			
				rate, _ := record.Get("rate")
			
				rateHost := domains.RateHost{
					Host: domains.Host{
						ID: hostID.(string),
						Email: hostEmail.(string),
						Username: hostUsername.(string),
					},
					Guest: domains.Guest{
						ID: guestID.(string),
						Email: guestEmail.(string),
						Username: guestUsername.(string),
					},
					Rate: rate.(int64),
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
