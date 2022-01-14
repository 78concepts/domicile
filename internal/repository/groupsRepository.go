package repository

import (
	"78concepts.com/domicile/internal/model"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type IGroupsRepository interface {
	GetGroups(ctx context.Context) ([]model.Group, error)
	GetGroup(ctx context.Context, id uint64) (*model.Group, error)
	CreateGroup(ctx context.Context, id uint64, name string) (*model.Group, error)
	UpdateGroup(ctx context.Context, id uint64, name string, active bool) (*model.Group, error)
	GetGroupMembers(ctx context.Context, id uint64) ([]model.GroupMember, error)
	CreateGroupMember(ctx context.Context, id uint64, ieeeAddress string) (*model.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id uint64, ieeeAddress string) error
}

type PostgresGroupsRepository struct {
	Postgres *pgxpool.Pool
}

func (r *PostgresGroupsRepository) GetGroups(ctx context.Context) ([]model.Group, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT * FROM GROUPS")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var objects []model.Group

	for rows.Next() {
		var row model.Group
		err = rows.Scan(&row.Id, &row.DateCreated, &row.DateModified, &row.FriendlyName, &row.Active)
		if err != nil {
			log.Fatal("GetGroups:", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetGroups:", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresGroupsRepository) GetGroup(ctx context.Context, id uint64) (*model.Group, error) {

	row := r.Postgres.QueryRow(ctx, "SELECT * FROM GROUPS WHERE ID = $1", id)

	var object model.Group

	err := row.Scan(&object.Id, &object.DateCreated, &object.DateModified, &object.FriendlyName, &object.Active)
	if err != nil {
		log.Fatal("GetGroup:", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresGroupsRepository) CreateGroup(ctx context.Context, id uint64, name string) (*model.Group, error) {

	dateCreated := time.Now().UTC()

	query := `
				INSERT INTO GROUPS
					(ID, DATE_CREATED, DATE_MODIFIED, FRIENDLY_NAME, ACTIVE)
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING ID, DATE_CREATED, DATE_MODIFIED, FRIENDLY_NAME, ACTIVE`

	row := r.Postgres.QueryRow(ctx, query, id, dateCreated, dateCreated, name, true)

	var object model.Group

	err := row.Scan(&object.Id, &object.DateCreated, &object.DateModified, &object.FriendlyName, &object.Active)

	if err != nil {
		log.Fatal("CreateGroup ", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresGroupsRepository) UpdateGroup(ctx context.Context, id uint64, name string, active bool) (*model.Group, error) {

	query := "UPDATE GROUPS SET FRIENDLY_NAME = $1, ACTIVE = $2 WHERE ID = $3 RETURNING ID, DATE_CREATED, DATE_MODIFIED, FRIENDLY_NAME, ACTIVE"

	row := r.Postgres.QueryRow(ctx, query, name, active, id)

	var object model.Group

	err := row.Scan(&object.Id, &object.DateCreated, &object.DateModified, &object.FriendlyName, &object.Active)

	if err != nil {
		log.Fatal("UpdateGroup", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresGroupsRepository) GetGroupMembers(ctx context.Context, id uint64) ([]model.GroupMember, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT GROUP_ID, IEEE_ADDRESS FROM GROUPS_DEVICES WHERE GROUP_ID = $1", id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var objects []model.GroupMember

	for rows.Next() {
		var row model.GroupMember
		err = rows.Scan(&row.GroupId, &row.IeeeAddress)
		if err != nil {
			log.Fatal("GetGroupMembers:", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetGroupMembers:", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresGroupsRepository) CreateGroupMember(ctx context.Context, id uint64, ieeeAddress string) (*model.GroupMember, error) {

	query := `
				INSERT INTO GROUPS_DEVICES
					(GROUP_ID, IEEE_ADDRESS)
				VALUES
					($1, $2)
				RETURNING GROUP_ID, IEEE_ADDRESS`

	row := r.Postgres.QueryRow(ctx, query, id, ieeeAddress)

	var object model.GroupMember

	err := row.Scan(&object.GroupId, &object.IeeeAddress)

	if err != nil {
		log.Fatal("CreateGroupMember ", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresGroupsRepository) DeleteGroupMember(ctx context.Context, id uint64, ieeeAddress string) (error) {

	query := "DELETE FROM GROUPS_DEVICES WHERE GROUP_ID = $1 AND IEEE_ADDRESS = $2"

	_, err := r.Postgres.Exec(ctx, query, id, ieeeAddress)

	if err != nil {
		log.Fatal("DeleteGroupMember ", err)
		return err
	}

	return nil
}

