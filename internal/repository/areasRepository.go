package repository

import (
	"78concepts.com/domicile/internal/model"
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type IAreasRepository interface {
	GetAreas(ctx context.Context) ([]model.Area, error)
	GetArea(ctx context.Context, uuid uuid.UUID) (*model.Area, error)
}

type PostgresAreasRepository struct {
	Postgres *pgxpool.Pool
}

func (r *PostgresAreasRepository) GetAreas(ctx context.Context) ([]model.Area, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT ID, UUID, DATE_CREATED, NAME FROM AREAS")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var objects []model.Area

	for rows.Next() {
		var row model.Area
		err = rows.Scan(&row.Id, &row.Uuid, &row.DateCreated, &row.Name)
		if err != nil {
			log.Fatal("GetAreas: ", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetAreas: ", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresAreasRepository) GetArea(ctx context.Context, uuid uuid.UUID) (*model.Area, error) {

	row := r.Postgres.QueryRow(ctx, "SELECT ID, UUID, DATE_CREATED, NAME FROM AREAS WHERE UUID = $1", uuid)

	var object model.Area

	err := row.Scan(&object.Id, &object.Uuid, &object.DateCreated, &object.Name)

	if err != nil {
		log.Fatal("GetArea: ", err)
		return nil, err
	}

	return &object, nil
}
