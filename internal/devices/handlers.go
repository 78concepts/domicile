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
	TopicGroups = broker.TopicRoot + "/bridge/groups"
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

func HandleGroups(mqttClient broker.MqttClient, devicesService *Service) {

	if token := mqttClient.Conn.Subscribe(TopicGroups, 0, func(client mqtt.Client, msg mqtt.Message) {
		groupsHandler(mqttClient.Ctx, devicesService, client, msg);
	}); token.Wait() && token.Error() != nil {
		log.Fatal("HandleGroups: Subscribe error: %s", token.Error())
		return
	}

}

var devicesHandler = func(ctx context.Context, devicesService *Service, client mqtt.Client, msg mqtt.Message) {

	log.Printf("Received devices message: %s from topic: %s\n", msg.Payload(), msg.Topic())

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

var groupsHandler = func(ctx context.Context, devicesService *Service, client mqtt.Client, msg mqtt.Message) {

	log.Printf("Received groups message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var objects []map[string]interface{}

	err := json.Unmarshal(msg.Payload(), &objects)

	if err != nil {
		log.Fatal(err)
	}

	groups, err:= devicesService.GetGroups(ctx)

	for _, object := range objects {

		var found *Group

		for i := range groups {
			if groups[i].Id == uint64(object["id"].(float64)) {
				found = &groups[i]
				break
			}
		}

		if found != nil {

			object["active"] = true

			if !found.Active {
				object["active"] = true
				devicesService.UpdateGroup(ctx, object)
			}

			if found.FriendlyName != object["friendly_name"] {
				devicesService.UpdateGroup(ctx, object)
			}

		} else {
			found, _ = devicesService.CreateGroup(ctx, object)
		}

		handleGroupMembers(ctx, devicesService, found, object["members"].([]interface{}))
	}

	// If a group in the database is no longer being reported, mark it as inactive
	for _, group := range groups {

		var found *map[string]interface {}

		for i := range objects {
			if uint64(objects[i]["id"].(float64)) == group.Id && group.Active {
				found = &objects[i]
				break
			}
		}

		if found == nil {

			object := map[string]interface{} {
				"id":  group.Id,
				"friendly_name":  group.FriendlyName,
				"active":  false,
			}

			devicesService.UpdateGroup(ctx, object)
		}
	}
}

var handleGroupMembers = func(ctx context.Context, devicesService *Service, group *Group, objects []interface{}) {

	if group == nil {
		return
	}

	groupMembers, _ := devicesService.GetGroupMembers(ctx, group.Id)
	for _, object := range objects {

		var found *GroupMember

		for i := range groupMembers {
			if groupMembers[i].IeeeAddress == object.(map[string]interface{})["ieee_address"].(string) {
				found = &groupMembers[i]
				break
			}
		}

		if found == nil {
			found, _ = devicesService.CreateGroupMember(ctx, group.Id, object.(map[string]interface{})["ieee_address"].(string))
		}
	}

	// If a group member in the database is no longer being reported, delete it
	for _, groupMember := range groupMembers {

		var found *map[string]interface {}

		for i := range objects {
			if objects[i].(map[string]interface{})["ieee_address"].(string) == groupMember.IeeeAddress {
				foundMember := objects[i].(map[string]interface{})
				found = &foundMember
				break
			}
		}

		if found == nil {
			devicesService.DeleteGroupMember(ctx, groupMember.GroupId, groupMember.IeeeAddress)
		}
	}
}

func TurnGroupOn(mqttClient broker.MqttClient, group *Group) {
	payload := "{\"state\": \"on\"}"

	if token := mqttClient.Conn.Publish(broker.TopicRoot + "/" + group.FriendlyName + "/set", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}


func TurnGroupOff(mqttClient broker.MqttClient, group *Group) {

	payload := "{\"state\": \"off\"}"

	if token := mqttClient.Conn.Publish(broker.TopicRoot + "/" + group.FriendlyName + "/set", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}