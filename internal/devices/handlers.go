package devices

import (
	"78concepts.com/domicile/internal/broker"
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

const (
	TopicGetDevices = broker.TopicRoot + "/bridge/config/devices/get"
	TopicDevices = broker.TopicRoot + "/bridge/devices"
)

func HandleDevices(mqttClient broker.MqttClient, devicesService *Service) {

	// Publish a message to trigger the broker to broadcast available devices
	if token := mqttClient.Conn.Publish(TopicGetDevices, 0, false, "{}"); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := mqttClient.Conn.Subscribe(TopicDevices, 0, func(client mqtt.Client, msg mqtt.Message) {
		devicesHandler(mqttClient.Ctx, devicesService, client, msg);
	}); token.Wait() && token.Error() != nil {
		log.Fatal("HandleDevices: Subscribe error: %s", token.Error())
		return
	}
}

var devicesHandler = func(ctx context.Context, devicesService *Service, client mqtt.Client, msg mqtt.Message) {

	//log.Printf("Received devices message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var objects []map[string]interface{}

	err := json.Unmarshal(msg.Payload(), &objects)

	if err != nil {
		log.Fatal(err)
	}

	devices, err:= devicesService.GetDevices(ctx)

	for _, object := range objects {

		var found *Device

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
				devicesService.UpdateDevice(ctx, object)
			}

			if found.FriendlyName != object["friendly_name"] {
				devicesService.UpdateDevice(ctx, object)
			}

		} else if object["type"] != "Coordinator" {
			found, _ = devicesService.CreateDevice(ctx, object)
		}

		if(found != nil) {
			if token := client.Subscribe(broker.TopicRoot+"/"+found.FriendlyName, 0, func(client mqtt.Client, msg mqtt.Message) {
				deviceHandler(ctx, devicesService, found, client, msg);
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

			devicesService.UpdateDevice(ctx, object)
		}
	}
}

var deviceHandler = func(ctx context.Context, devicesService *Service, device *Device, client mqtt.Client, msg mqtt.Message) {

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
		devicesService.UpdateDeviceBattery(ctx, map[string]interface{}{
			"ieee_address": device.IeeeAddress,
			"battery":      report["battery"],
		})
	}

	if device.AreaId == nil {
		log.Println("DeviceHandler: device does not belong to an area, will not save report")
		return
	}

	//TODO
	// 2022/01/11 13:26:35 map[battery:91 humidity:65.79 linkquality:255 pressure:1003.1 temperature:27.21 voltage:2985]
	// Groups

	if report["temperature"] != nil && device.ModelId == "lumi.weather" { // Celcius
		devicesService.CreateTemperatureReport(ctx, device.IeeeAddress, *device.AreaId, report["temperature"].(float64))
	}
	if report["humidity"] != nil { // Percentage
		devicesService.CreateHumidityReport(ctx, device.IeeeAddress, *device.AreaId, report["humidity"].(float64))
	}
	if report["pressure"] != nil { // HectoPascals
		devicesService.CreatePressureReport(ctx, device.IeeeAddress, *device.AreaId, report["pressure"].(float64))
	}
	if report["illuminance"] != nil && report["illuminance_lux"] != nil { // Raw illuminance / Lux
		devicesService.CreateIlluminanceReport(ctx, device.IeeeAddress, *device.AreaId, report["illuminance"].(float64), report["illuminance_lux"].(float64))
	}
}

//	// Need to get all of the devices from the DB
//	// IF the device is not there, then update the DB and set to active
//	// Update the friendly name base on the ieee in the list

//	//channel := make(chan int)
//	//go subscribeToGroups(channel, client)
//}
//
//func subscribeToSensorTopics(client mqtt.Client, devices []Device) {
//
//	for _, device := range devices {
//		if device.Type == "EndDevice" {
//			switch device.ModelId {
//				case "lumi.weather":
//					channel := make(chan int)
//					go subscribeToSensorTopic(channel, client, device)
//					break
//				default:
//					break
//			}
//		}
//	}
//}
//
//func subscribeToSensorTopic(c chan <- int, client mqtt.Client, device Device) {
//
//	defer close(c)
//
//	if token := client.Subscribe(TopicRoot + "/" + device.FriendlyName, 0, sensorSubscribeHandler); token.Wait() && token.Error() != nil {
//		panic(token.Error())
//	}
//}

