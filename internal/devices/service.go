package devices

import (
	"context"
	"errors"
	"time"
)

func NewService(devicesRepository IRepository) *Service {
	return &Service{devicesRepository: devicesRepository}
}

type Service struct {
	devicesRepository IRepository
}

func (s *Service) GetDevices(ctx context.Context) (devices []Device, err error) {
	return s.devicesRepository.GetDevices(ctx)
}

func (s *Service) CreateDevice(ctx context.Context, object map[string]interface{}) (device *Device, err error) {

	if object == nil {
		return nil, errors.New("CreateDevice: Object is null")
	}

	var lastSeen *uint64;
	if object["last_seen"] != nil {
		x := object["last_seen"].(uint64)
		lastSeen = &x
	}

	var description *string;
	if object["description"] != nil {
		x := object["description"].(string)
		description = &x
	}

	var model *string;
	if object["model"] != nil {
		x := object["model"].(string)
		model = &x
	}

	var vendor *string;
	if object["vendor"] != nil {
		x := object["vendor"].(string)
		vendor = &x
	}

	return s.devicesRepository.CreateDevice(
		ctx,
		object["ieee_address"].(string),
		object["date_code"].(string),
		object["friendly_name"].(string),
		description,
		object["manufacturer"].(string),
		model,
		object["model_id"].(string),
		lastSeen,
		vendor,
		object["type"].(string),
	)
}

func (s *Service) UpdateDevice(ctx context.Context, object map[string]interface{}) (device *Device, err error) {

	if object == nil {
		return nil, errors.New("UpdateDevice: Object is null")
	}

	if object["ieee_address"] == nil {
		return nil, errors.New("UpdateDevice: IEEE address is null")
	}

	return s.devicesRepository.UpdateDevice(
		ctx,
		object["ieee_address"].(string),
		object["friendly_name"].(string),
		object["active"].(bool),
	)
}

func (s *Service) UpdateDeviceBattery(ctx context.Context, object map[string]interface{}) (device *Device, err error) {

	if object == nil {
		return nil, errors.New("UpdateDevice: Object is null")
	}

	if object["ieee_address"] == nil {
		return nil, errors.New("UpdateDevice: IEEE address is null")
	}

	if object["battery"] != nil {
		return s.devicesRepository.UpdateDeviceBattery(
			ctx,
			object["ieee_address"].(string),
			object["battery"].(float64),
		)
	} else {
		return nil, errors.New("Battery not found")
	}
}

func (s *Service) CreateTemperatureReport(ctx context.Context, deviceId string, areaId uint64, value float64) (device *TemperatureReport, err error) {

	return s.devicesRepository.CreateTemperatureReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *Service) CreateHumidityReport(ctx context.Context, deviceId string, areaId uint64, value float64) (device *HumidityReport, err error) {

	return s.devicesRepository.CreateHumidityReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *Service) CreatePressureReport(ctx context.Context, deviceId string, areaId uint64, value float64) (device *PressureReport, err error) {

	return s.devicesRepository.CreatePressureReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
	)
}

func (s *Service) CreateIlluminanceReport(ctx context.Context, deviceId string, areaId uint64, value float64, valueLux float64) (device *IlluminanceReport, err error) {

	return s.devicesRepository.CreateIlluminanceReport(
		ctx,
		deviceId,
		areaId,
		time.Now().UTC(),
		value,
		valueLux,
	)
}