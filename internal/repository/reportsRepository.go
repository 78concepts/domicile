package repository

import (
	"78concepts.com/domicile/internal/model"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type IReportsRepository interface {
	CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.TemperatureReport, error)
	CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.HumidityReport, error)
	CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.PressureReport, error)
	CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64, valueLux float64) (*model.IlluminanceReport, error)
	GetTemperatureReports(ctx context.Context, areaId uint64, startDate time.Time, endDate time.Time) ([]model.TemperatureReport, error)
	GetHumidityReports(ctx context.Context, areaId uint64) ([]model.HumidityReport, error)
	GetPressureReports(ctx context.Context, areaId uint64) ([]model.PressureReport, error)
	GetIlluminanceReports(ctx context.Context, areaId uint64) ([]model.IlluminanceReport, error)
}

type PostgresReportsRepository struct {
	Postgres *pgxpool.Pool
}

func (r *PostgresReportsRepository) CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.TemperatureReport, error) {

	query := "INSERT INTO TEMPERATURE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object model.TemperatureReport

	err := row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreateTemperatureReport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresReportsRepository) CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.HumidityReport, error) {

	query := "INSERT INTO HUMIDITY_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object model.HumidityReport

	err := row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreateHumidityReport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresReportsRepository) CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64) (*model.PressureReport, error) {

	query := "INSERT INTO PRESSURE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE) VALUES ($1, $2, $3, $4) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value)

	var object model.PressureReport

	err := row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value)

	if err != nil {
		log.Fatal("CreatePressureRport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresReportsRepository) CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, date time.Time, value float64, valueLux float64) (*model.IlluminanceReport, error) {

	query := "INSERT INTO ILLUMINANCE_REPORTS (DEVICE_ID, AREA_ID, DATE, VALUE, VALUE_LUX) VALUES ($1, $2, $3, $4, $5) RETURNING DEVICE_ID, AREA_ID, DATE, VALUE, VALUE_LUX"

	row := r.Postgres.QueryRow(ctx, query, deviceId, areaId, date, value, valueLux)

	var object model.IlluminanceReport

	err := row.Scan(&object.DeviceId, &object.AreaId, &object.Date, &object.Value, &object.ValueLux)

	if err != nil {
		log.Fatal("CreateTemperatureReport", err)
		return nil, err
	}

	return &object, nil
}

func (r *PostgresReportsRepository) GetTemperatureReports(ctx context.Context, areaId uint64, startDate time.Time, endDate time.Time) ([]model.TemperatureReport, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT DEVICE_ID, AREA_ID, DATE, VALUE FROM TEMPERATURE_REPORTS WHERE AREA_ID = $1 AND DATE >= $2 AND DATE <= $3 ORDER BY DATE ASC", areaId, startDate, endDate)
log.Println(startDate)
	log.Println(endDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]model.TemperatureReport, 0)

	for rows.Next() {
		var row model.TemperatureReport
		err = rows.Scan(&row.DeviceId, &row.AreaId, &row.Date, &row.Value)
		if err != nil {
			log.Fatal("GetTemperatureReports:", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetTemperatureReports:", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresReportsRepository) GetHumidityReports(ctx context.Context, areaId uint64) ([]model.HumidityReport, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT DEVICE_ID, AREA_ID, DATE, VALUE FROM HUMIDITY_REPORTS WHERE AREA_ID=$1 ORDER BY DATE ASC", areaId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]model.HumidityReport, 0)

	for rows.Next() {
		var row model.HumidityReport
		err = rows.Scan(&row.DeviceId, &row.AreaId, &row.Date, &row.Value)
		if err != nil {
			log.Fatal("GetHumidityReports:", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetHumidityReports:", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresReportsRepository) GetPressureReports(ctx context.Context, areaId uint64) ([]model.PressureReport, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT DEVICE_ID, AREA_ID, DATE, VALUE FROM PRESSURE_REPORTS WHERE AREA_ID=$1 ORDER BY DATE ASC", areaId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]model.PressureReport, 0)

	for rows.Next() {
		var row model.PressureReport
		err = rows.Scan(&row.DeviceId, &row.AreaId, &row.Date, &row.Value)
		if err != nil {
			log.Fatal("GetPressureReports:", err)
			return nil, err
		}

		objects = append(objects, row)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("GetPressureReports:", err)
		return nil, err
	}

	return objects, nil
}

func (r *PostgresReportsRepository) GetIlluminanceReports(ctx context.Context, areaId uint64) ([]model.IlluminanceReport, error) {

	rows, err := r.Postgres.Query(ctx, "SELECT DEVICE_ID, AREA_ID, DATE, VALUE FROM ILLUMINANCE_REPORTS WHERE AREA_ID=$1 ORDER BY DATE ASC", areaId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := make([]model.IlluminanceReport, 0)

	for rows.Next() {
		var row model.IlluminanceReport
		err = rows.Scan(&row.DeviceId, &row.AreaId, &row.Date, &row.Value)
		if err != nil {
			log.Fatal("GetIlluminanceReports:", err)
			return nil, err
		}

		objects = append(objects, row)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("GetIlluminanceReports:", err)
		return nil, err
	}

	return objects, nil
}