package devices

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"log"
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

	var dateCode *string;
	if object["date_code"] != nil {
		x := object["date_code"].(string)
		dateCode = &x
	}

	var manufacturer *string;
	if object["manufacturer"] != nil {
		x := object["manufacturer"].(string)
		manufacturer = &x
	}

	var modelId *string;
	if object["model_id"] != nil {
		x := object["model_id"].(string)
		modelId = &x
	}

	var lastSeen *uint64;
	if object["last_seen"] != nil {
		x := object["last_seen"].(uint64)
		lastSeen = &x
	}

	var deviceType *string;
	if object["type"] != nil {
		x := object["type"].(string)
		deviceType = &x
	}

	log.Println(object)

	return s.devicesRepository.CreateDevice(
		ctx,
		object["ieee_address"].(string),
		dateCode,
		object["friendly_name"].(string),
		manufacturer,
		modelId,
		lastSeen,
		deviceType,
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

func (s *Service) GetAreas(ctx context.Context) (areas []Area, err error) {
	return s.devicesRepository.GetAreas(ctx)
}

func (s *Service) GetArea(ctx context.Context, uuid uuid.UUID) (areas *Area, err error) {
	return s.devicesRepository.GetArea(ctx, uuid)
}

func (s *Service) GetTemperatureReports(ctx context.Context, areaId uint64) (reports []TemperatureReport, err error) {
	return s.devicesRepository.GetTemperatureReports(ctx, areaId)
}

func (s *Service) GetHumidityReports(ctx context.Context, areaId uint64) (reports []HumidityReport, err error) {
	return s.devicesRepository.GetHumidityReports(ctx, areaId)
}

func (s *Service) GetGroups(ctx context.Context) (groups []Group, err error) {
	return s.devicesRepository.GetGroups(ctx)
}

func (s *Service) GetGroup(ctx context.Context, id uint64) (group *Group, err error) {
	return s.devicesRepository.GetGroup(ctx, id)
}

func (s *Service) CreateGroup(ctx context.Context, object map[string]interface{}) (device *Group, err error) {

	if object == nil {
		return nil, errors.New("CreateGroup: Object is null")
	}

	return s.devicesRepository.CreateGroup(
		ctx,
		uint64(object["id"].(float64)),
		object["friendly_name"].(string),
	)
}

func (s *Service) UpdateGroup(ctx context.Context, object map[string]interface{}) (group *Group, err error) {

	if object == nil {
		return nil, errors.New("UpdateGroup: Object is null")
	}

	return s.devicesRepository.UpdateGroup(
		ctx,
		uint64(object["id"].(float64)),
		object["friendly_name"].(string),
		object["active"].(bool),
	)
}

func (s *Service) GetGroupMembers(ctx context.Context, id uint64) (members []GroupMember, err error) {
	return s.devicesRepository.GetGroupMembers(ctx, id)
}

func (s *Service) CreateGroupMember(ctx context.Context, id uint64, ieeeAddress string) (groupMember *GroupMember, err error) {
	return s.devicesRepository.CreateGroupMember(ctx, id, ieeeAddress)
}

func (s *Service) DeleteGroupMember(ctx context.Context, id uint64, ieeeAddress string) (err error) {
	return s.devicesRepository.DeleteGroupMember(ctx, id, ieeeAddress)
}

