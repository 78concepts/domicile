package devices

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type IRepository interface {
	GetDevices(ctx context.Context) ([]Device, error)
	CreateDevice(ctx context.Context, ieeeAddress string, dateCode string, name string, description *string, manufacturer string, model *string, modelId string, lastSeen *uint64, vendor *string, deviceType string) (*Device, error)
	UpdateDevice(ctx context.Context, ieeeAddress string, name string, active bool) (*Device, error)
	UpdateDeviceBattery(ctx context.Context, ieeeAddress string, battery float64) (*Device, error)
	CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*TemperatureReport, error)
	CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*HumidityReport, error)
	CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*PressureReport, error)
	CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64, valueLux float64) (*IlluminanceReport, error)
}

type PostgresRepository struct {
	Postgres *pgxpool.Pool
}

var deviceReturnFields = "IEEE_ADDRESS, DATE_CREATED, DATE_MODIFIED, DATE_CODE, FRIENDLY_NAME, AREA_ID, DESCRIPTION, MANUFACTURER, MODEL, MODEL_ID, LAST_SEEN, VENDOR, TYPE, BATTERY, ACTIVE"

var scanDeviceRows = func(rows pgx.Rows, row *Device) error {
	return rows.Scan(&row.IeeeAddress, &row.DateCreated, &row.DateModified, &row.DateCode, &row.FriendlyName, &row.AreaId, &row.Description, &row.Manufacturer, &row.Model, &row.ModelId, &row.LastSeen, &row.Vendor, &row.Type, &row.Battery, &row.Active)
}

var scanDeviceRow = func(rows pgx.Row, row *Device) error {
	return rows.Scan(&row.IeeeAddress, &row.DateCreated, &row.DateModified, &row.DateCode, &row.FriendlyName, &row.AreaId, &row.Description, &row.Manufacturer, &row.Model, &row.ModelId, &row.LastSeen, &row.Vendor, &row.Type, &row.Battery, &row.Active)
}

func (r *PostgresRepository) GetDevices(ctx context.Context) (result []Device, err error) {

	rows, err := r.Postgres.Query(ctx, "SELECT * FROM DEVICES")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var objects []Device

	for rows.Next() {
		var row Device
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

func (r *PostgresRepository) CreateDevice(ctx context.Context, ieeeAddress string, dateCode string, name string, description *string, manufacturer string, model *string, modelId string, lastSeen *uint64, vendor *string, deviceType string) (result *Device, err error) {

	dateCreated := time.Now().UTC()

	query := `
				INSERT INTO DEVICES 
					(IEEE_ADDRESS, DATE_CREATED, DATE_MODIFIED, DATE_CODE, FRIENDLY_NAME, 
					 DESCRIPTION, MANUFACTURER, MODEL, MODEL_ID, LAST_SEEN, VENDOR, TYPE, ACTIVE) 
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
				RETURNING ` + deviceReturnFields

	row := r.Postgres.QueryRow(ctx, query, ieeeAddress, dateCreated, dateCreated, dateCode, name,
		description, manufacturer, model, modelId, lastSeen, vendor, deviceType, true)

	var object Device

	err = scanDeviceRow(row, &object)

	if err != nil {
		log.Fatal("CreateDevice ", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) UpdateDevice(ctx context.Context, ieeeAddress string, name string, active bool) (result *Device, err error) {

	query := "UPDATE DEVICES SET FRIENDLY_NAME = $1, ACTIVE = $2 WHERE IEEE_ADDRESS = $3 RETURNING " + deviceReturnFields

	row := r.Postgres.QueryRow(ctx, query, name, active, ieeeAddress)

	var object Device

	err = scanDeviceRow(row, &object)

	if err != nil {
		log.Fatal("UpdateDevice", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) UpdateDeviceBattery(ctx context.Context, ieeeAddress string, battery float64) (result *Device, err error) {

	dateModified := time.Now().UTC()

	query := "UPDATE DEVICES SET DATE_MODIFIED = $1, LAST_SEEN = $2, BATTERY = $3 WHERE IEEE_ADDRESS = $4 RETURNING " + deviceReturnFields

	row := r.Postgres.QueryRow(ctx, query, dateModified, dateModified, battery, ieeeAddress)

	var object Device

	err = scanDeviceRow(row, &object)

	if err != nil {
		log.Fatal("UpdateDeviceBattery", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (result *TemperatureReport, err error) {

	query := "INSERT INTO TEMPERATURE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object TemperatureReport

	err = row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreateTemperatureReport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (result *HumidityReport, err error) {

	query := "INSERT INTO HUMIDITY_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object HumidityReport

	err = row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreateHumidityReport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (result *PressureReport, err error) {

	query := "INSERT INTO PRESSURE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object PressureReport

	err = row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreatePressureRport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresRepository) CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64, valueLux float64) (result *IlluminanceReport, err error) {

	query := "INSERT INTO ILLUMINANCE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE, VALUE_LUX) VALUES ($1, $2, $3, $4, $5) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE, VALUE_LUX"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value, valueLux)

	var object IlluminanceReport

	err = row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value, &object.ValueLux)

	if err != nil {
		log.Fatal("CreateTemperatureReport", err)
		return nil, err
	}

	return &object, nil
}