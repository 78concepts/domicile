package service

import (
	"78concepts.com/domicile/internal/model"
	"78concepts.com/domicile/internal/repository"
	"context"
	"time"
)

func NewReportsService(reportsRepository repository.IReportsRepository) *ReportsService {
	return &ReportsService{reportsRepository: reportsRepository}
}

type ReportsService struct {
	reportsRepository repository.IReportsRepository
}

func (s *ReportsService) CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, value float64) (*model.TemperatureReport, error) {

	return s.reportsRepository.CreateTemperatureReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *ReportsService) CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, value float64) (*model.HumidityReport, error) {

	return s.reportsRepository.CreateHumidityReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *ReportsService) CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, value float64) (*model.PressureReport, error) {

	return s.reportsRepository.CreatePressureReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *ReportsService) CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, value float64, valueLux float64) (*model.IlluminanceReport, error) {

	return s.reportsRepository.CreateIlluminanceReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
		valueLux,
	)
}

func (s *ReportsService) GetTemperatureReports(ctx context.Context, areaId uint64, startDate time.Time, endDate time.Time) ([]model.TemperatureReport, error) {
	return s.reportsRepository.GetTemperatureReports(ctx, areaId, startDate, endDate)
}

func (s *ReportsService) GetHumidityReports(ctx context.Context, areaId uint64) ([]model.HumidityReport, error) {
	return s.reportsRepository.GetHumidityReports(ctx, areaId)
}

func (s *ReportsService) GetPressureReports(ctx context.Context, areaId uint64) ([]model.PressureReport, error) {
	return s.reportsRepository.GetPressureReports(ctx, areaId)
}

func (s *ReportsService) GetIlluminanceReports(ctx context.Context, areaId uint64) ([]model.IlluminanceReport, error) {
	return s.reportsRepository.GetIlluminanceReports(ctx, areaId)
}

