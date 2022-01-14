package api

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
)

func NewDevicesApi(ctx context.Context, client *broker.MqttClient, devicesService *service.DevicesService) *DevicesApi {
	return &DevicesApi{ctx: ctx, client: client, devicesService: devicesService}
}

type DevicesApi struct {
	ctx context.Context
	client *broker.MqttClient
	devicesService *service.DevicesService
}

func (a *DevicesApi) GetState(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, string("Sd"))
	//s := "zigbee2mqtt/Office Ceiling Light 1/get"
	log.Println(a.devicesService.RequestDeviceState(a.client, "Office Ceiling Light 1"))

	//fmt.Fprintf(w, service.DevicesService.GetDeviceState(a.client, "Office Ceiling Light 1"))

}