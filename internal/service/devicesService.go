package service

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/model"
	"78concepts.com/domicile/internal/repository"
	"context"
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

const (
	TopicGetDevices = broker.TopicRoot + "/bridge/config/devices/get"
	TopicDevices = broker.TopicRoot + "/bridge/devices"
)

func NewDevicesService(reportsService *ReportsService, devicesRepository repository.IDevicesRepository) *DevicesService {
	return &DevicesService{reportsService: reportsService, devicesRepository: devicesRepository}
}

type DevicesService struct {
	reportsService *ReportsService
	devicesRepository repository.IDevicesRepository
}

func (s *DevicesService) ManageDevices(mqttClient *broker.MqttClient) {

	// Publish a message to trigger the broker to broadcast available devices
	if token := mqttClient.Conn.Publish(TopicGetDevices, 0, false, "{}"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := mqttClient.Conn.Subscribe(TopicDevices, 0, func(client mqtt.Client, msg mqtt.Message) {
		s.HandleDevicesMessage(mqttClient.Ctx, client, msg);
	}); token.Wait() && token.Error() != nil {
		log.Fatal("HandleDevices: Subscribe error: %s", token.Error())
		return
	}
}

func (s *DevicesService) HandleDevicesMessage(ctx context.Context, client mqtt.Client, msg mqtt.Message) {

	log.Printf("Received devices message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var objects []map[string]interface{}

	err := json.Unmarshal(msg.Payload(), &objects)

	if err != nil {
		log.Fatal(err)
	}

	devices, err:= s.GetDevices(ctx)

	for _, object := range objects {

		var found *model.Device

		for i := range devices {
			if devices[i].IeeeAddress == object["ieee_address"] {
				found = &devices[i]
				break
			}
		}

		if found != nil {

			object["active"] = true

			if !found.Active {
				object["active"] = true
				s.UpdateDevice(ctx, object)
			}

			if found.FriendlyName != object["friendly_name"] {
				s.UpdateDevice(ctx, object)
			}

		} else if object["type"] != "Coordinator" {
			found, _ = s.CreateDevice(ctx, object)
		}

		if found != nil {
			if token := client.Subscribe(broker.TopicRoot+"/"+found.FriendlyName, 0, func(client mqtt.Client, msg mqtt.Message) {
				s.HandleDeviceMessage(ctx, msg, found);
			}); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}

	}

	// If a device in the database is no longer being reported, mark it as inactive
	for _, device := range devices {

		var found *map[string]interface {}

		for i := range objects {
			if objects[i]["ieee_address"] == device.IeeeAddress && device.Active {
				found = &objects[i]
				break
			}
		}

		if found == nil {

			object := map[string]interface{} {
				"ieee_address":  device.IeeeAddress,
				"friendly_name":  device.FriendlyName,
				"active":  false,
			}

			s.UpdateDevice(ctx, object)
		}
	}
}

func (s *DevicesService) HandleDeviceMessage(ctx context.Context, msg mqtt.Message, device *model.Device) {

	log.Printf("Received device message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	if len(msg.Payload()) == 0 {
		return
	}

	var report map[string]interface{}

	err := json.Unmarshal(msg.Payload(), &report)

	if err != nil {
		log.Fatal("DeviceHandler", err)
		log.Fatal(err)
	}

	if report["battery"] != nil {
		s.UpdateDeviceBattery(ctx, map[string]interface{}{
			"ieee_address": device.IeeeAddress,
			"battery":      report["battery"],
		})
	}

	if device.AreaId == nil {
		log.Println("DeviceHandler: device does not belong to an area, will not save report")
		return
	}

	if report["temperature"] != nil && device.ModelId != nil && *device.ModelId == "lumi.weather" { // Celcius
		s.reportsService.CreateTemperatureReport(ctx, device.IeeeAddress, *device.AreaId, report["temperature"].(float64))
	}
	if report["humidity"] != nil { // Percentage
		s.reportsService.CreateHumidityReport(ctx, device.IeeeAddress, *device.AreaId, report["humidity"].(float64))
	}
	if report["pressure"] != nil { // HectoPascals
		s.reportsService.CreatePressureReport(ctx, device.IeeeAddress, *device.AreaId, report["pressure"].(float64))
	}
	if report["illuminance"] != nil && report["illuminance_lux"] != nil { // Raw illuminance / Lux
		s.reportsService.CreateIlluminanceReport(ctx, device.IeeeAddress, *device.AreaId, report["illuminance"].(float64), report["illuminance_lux"].(float64))
	}
}

func (s *DevicesService) GetDevices(ctx context.Context) ([]model.Device, error) {
	return s.devicesRepository.GetDevices(ctx)
}

func (s *DevicesService) CreateDevice(ctx context.Context, object map[string]interface{}) (*model.Device, error) {

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

func (s *DevicesService) UpdateDevice(ctx context.Context, object map[string]interface{}) (*model.Device, error) {

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

func (s *DevicesService) UpdateDeviceBattery(ctx context.Context, object map[string]interface{}) (*model.Device, error) {

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
