package repository

import "htmx/entity"

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
