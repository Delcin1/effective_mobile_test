package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

/*
{
  "regNum": "X123XX150",
  "mark": "Lada",
  "model": "Vesta",
  "year": 2002,
  "owner": {
    "name": "string",
    "surname": "string",
    "patronymic": "string"
  }
}
*/

// @Schema
type Owner struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

// @Schema
type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  Owner  `json:"owner"`
}

// @Schema
type SearchRequest struct {
	Query    string `json:"query"`
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
}

type Storage struct {
	db *sql.DB
}

func New(dbUrl string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveOwner(owner Owner) (int, error) {
	const op = "storage.postgres.SaveOwner"

	var id int
	err := s.db.QueryRow("INSERT INTO owners(name, surname, patronymic) VALUES ($1, $2, $3) RETURNING owner_id",
		owner.Name, owner.Surname, owner.Patronymic).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetOwnerID(owner Owner) (int, error) {
	const op = "storage.postgres.GetOwnerID"

	rows, err := s.db.Query("SELECT id FROM owners WHERE name = $1 AND surname = $2 AND patronymic = $3",
		owner.Name, owner.Surname, owner.Patronymic)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		id, err = s.SaveOwner(owner)
		if err != nil {
			return -1, fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) SaveCar(car Car) (int, error) {
	const op = "storage.postgres.SaveCar"

	var id int
	err := s.db.QueryRow("INSERT INTO cars(reg_num, mark, model, year) VALUES ($1, $2, $3, $4) RETURNING car_id",
		car.RegNum, car.Mark, car.Model, car.Year).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	ownerID, err := s.GetOwnerID(car.Owner)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("INSERT INTO cars_owners(car_reg_num, owner_id) VALUES ($1, $2)",
		car.RegNum, ownerID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetCarsBySearchRequest(searchRequest SearchRequest) ([]Car, error) {
	const op = "storage.postgres.GetCarsBySearchRequest"

	var cars []Car

	rows, err := s.db.Query(`SELECT DISTINCT c.reg_num, c.mark, c.model, c.year, o.name, o.surname, o.patronymic 
								   FROM cars c 
    							   JOIN cars_owners co ON c.car_id = co.car_id 
								   JOIN owners o ON co.owner_id = o.owner_id
								   WHERE c.reg_num LIKE $1 or c.mark LIKE $1 or c.model LIKE $1 or o.name LIKE $1
								      or o.surname LIKE $1 or o.patronymic LIKE $1
								   LIMIT $2 OFFSET $2*($3-1)`, "%"+searchRequest.Query+"%", searchRequest.PageSize, searchRequest.PageNum)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var car Car
		var owner Owner
		err = rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &owner.Name, &owner.Surname, &owner.Patronymic)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		car.Owner = owner
		cars = append(cars, car)
	}

	return cars, nil
}

func (s *Storage) DeleteCar(carID int) error {
	const op = "storage.postgres.DeleteCar"

	_, err := s.db.Exec("DELETE FROM cars_owners WHERE car_id = $1", carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM cars WHERE car_id = $1", carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteOwner(ownerID int) error {
	const op = "storage.postgres.DeleteOwner"

	_, err := s.db.Exec("DELETE FROM cars_owners WHERE owner_id = $1", ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM owners WHERE owner_id = $1", ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *Storage) UpdateRegNum(carID int, newRegNum string) error {
	const op = "storage.postgres.UpdateRegNum"

	_, err := s.db.Exec("UPDATE cars SET reg_num = $1 WHERE car_id = $2", newRegNum, carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateMark(carID int, newMark string) error {
	const op = "storage.postgres.UpdateMark"

	_, err := s.db.Exec("UPDATE cars SET mark = $1 WHERE car_id = $2", newMark, carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateModel(carID int, newModel string) error {
	const op = "storage.postgres.UpdateModel"

	_, err := s.db.Exec("UPDATE cars SET model = $1 WHERE car_id = $2", newModel, carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateYear(carID int, newYear int) error {
	const op = "storage.postgres.UpdateYear"

	_, err := s.db.Exec("UPDATE cars SET year = $1 WHERE car_id = $2", newYear, carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateOwner(carID int, newOwner Owner) error {
	const op = "storage.postgres.UpdateOwner"

	ownerID, err := s.GetOwnerID(newOwner)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("DELETE FROM cars_owners WHERE car_id = $1", carID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec("INSERT INTO cars_owners(car_id, owner_id) VALUES ($1, $2)", carID, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateOwnerName(ownerID int, newName string) error {
	const op = "storage.postgres.UpdateOwnerName"

	_, err := s.db.Exec("UPDATE owners SET name = $1 WHERE owner_id = $2", newName, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateOwnerSurname(ownerID int, newSurname string) error {
	const op = "storage.postgres.UpdateOwnerSurname"

	_, err := s.db.Exec("UPDATE owners SET surname = $1 WHERE owner_id = $2", newSurname, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateOwnerPatronymic(ownerID int, newPatronymic string) error {
	const op = "storage.postgres.UpdateOwnerPatronymic"

	_, err := s.db.Exec("UPDATE owners SET patronymic = $1 WHERE owner_id = $2", newPatronymic, ownerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
