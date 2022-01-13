package repository

import (
	"78concepts.com/domicile/internal/model"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type IDevicesRepository interface {
	GetDevices(ctx context.Context) ([]model.Device, error)
	CreateDevice(ctx context.Context, ieeeAddress string, dateCode *string, name string, manufacturer *string, modelId *string, lastSeen *uint64, deviceType *string) (*model.Device, error)
	UpdateDevice(ctx context.Context, ieeeAddress string, name string, active bool) (*model.Device, error)
	UpdateDeviceBattery(ctx context.Context, ieeeAddress string, battery float64) (*model.Device, error)
}

type PostgresDevicesRepository struct {
	Postgres *pgxpool.Pool
}

var returnFields = "IEEE_ADDRESS, DATE_CREATED, DATE_MODIFIED, DATE_CODE, FRIENDLY_NAME, AREA_ID, MANUFACTURER, MODEL_ID, LAST_SEEN, TYPE, BATTERY, ACTIVE"

var scanDeviceRows = func(rows pgx.Rows, object *model.Device) error {
	return rows.Scan(&object.IeeeAddress, &object.DateCreated, &object.DateModified, &object.DateCode, &object.FriendlyName, &object.AreaId, &object.Manufacturer, &object.ModelId, &object.LastSeen, &object.Type, &object.Battery, &object.Active)
}

var scanRow = func(row pgx.Row, object *model.Device) error {
	return row.Scan(&object.IeeeAddress, &object.DateCreated, &object.DateModified, &object.DateCode, &object.FriendlyName, &object.AreaId, &object.Manufacturer, &object.ModelId, &object.LastSeen, &object.Type, &object.Battery, &object.Active)
}

func (r *PostgresDevicesRepository) GetDevices(ctx context.Context) (result []model.Device, err error) {

	rows, err := r.Postgres.Query(ctx, "SELECT * FROM DEVICES")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var objects []model.Device

	for rows.Next() {
		var row model.Device

		//TODO
		err = scanDeviceRows(rows, &row)
		if err != nil {
			log.Fatal("GetDevices: ", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetDevices: ", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresDevicesRepository) CreateDevice(ctx context.Context, ieeeAddress string, dateCode *string, name string, manufacturer *string, modelId *string, lastSeen *uint64, deviceType *string) (*model.Device, error) {

	dateCreated := time.Now().UTC()

	query := `
				INSERT INTO DEVICES 
					(IEEE_ADDRESS, DATE_CREATED, DATE_MODIFIED, DATE_CODE, FRIENDLY_NAME, 
					 MANUFACTURER, MODEL_ID, LAST_SEEN, TYPE, ACTIVE) 
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
				RETURNING ` + returnFields

	row := r.Postgres.QueryRow(ctx, query, ieeeAddress, dateCreated, dateCreated, dateCode, name,
		manufacturer, modelId, lastSeen, deviceType, true)

	var object model.Device

	err := scanRow(row, &object)

	if err != nil {
		log.Fatal("CreateDevice ", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresDevicesRepository) UpdateDevice(ctx context.Context, ieeeAddress string, name string, active bool) (*model.Device, error) {

	query := "UPDATE DEVICES SET FRIENDLY_NAME = $1, ACTIVE = $2 WHERE IEEE_ADDRESS = $3 RETURNING " + returnFields

	row := r.Postgres.QueryRow(ctx, query, name, active, ieeeAddress)

	var object model.Device

	err := scanRow(row, &object)

	if err != nil {
		log.Fatal("UpdateDevice", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresDevicesRepository) UpdateDeviceBattery(ctx context.Context, ieeeAddress string, battery float64) (*model.Device, error) {

	dateModified := time.Now().UTC()

	query := "UPDATE DEVICES SET DATE_MODIFIED = $1, LAST_SEEN = $2, BATTERY = $3 WHERE IEEE_ADDRESS = $4 RETURNING " + returnFields

	row := r.Postgres.QueryRow(ctx, query, dateModified, dateModified, battery, ieeeAddress)

	var object model.Device

	err := scanRow(row, &object)

	if err != nil {
		log.Fatal("UpdateDeviceBattery", err)
		return nil, err
	}

	return &object, nil
}
