package repository

import (
	"htmx/entity"

	"github.com/google/uuid"
)

func FindDogs() ([]entity.Dog, error) {
	results, err := DB.Query("SELECT * FROM dog")
	if err != nil {
		return nil, err
	}

	var dogs []entity.Dog

	for results.Next() {
		var dog entity.Dog

		err = results.Scan(&dog.ID, &dog.Name, &dog.Score, &dog.DateCreation)
		if err != nil {
			return nil, err
		}
		dogs = append(dogs, dog)
	}

	return dogs, nil
}

func IncrementScore(id uuid.UUID) (int, error) {
	stmt, err := DB.Prepare("UPDATE dog SET score = score + 1 WHERE id = ?")
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	var updatedScore int
	err = DB.QueryRow("SELECT score FROM dog WHERE id = ?", id).Scan(&updatedScore)
	if err != nil {
		return 0, err
	}

	return updatedScore, nil
}
