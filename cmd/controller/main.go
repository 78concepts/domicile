package main

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/database"
	"78concepts.com/domicile/internal/repository"
	"78concepts.com/domicile/internal/service"
	"context"
)

func main() {

	// Connect to the MQTT broker
	ctx, ctxCancel:= context.WithCancel(context.Background())

	var mqttClient = broker.NewMqttClient(ctx, ctxCancel, "api")

	// Connect to the database
	dbPool:= database.NewPGXPool()

	reportsService:= service.NewReportsService(&repository.PostgresReportsRepository{Postgres: dbPool})
	devicesService:= service.NewDevicesService(reportsService, &repository.PostgresDevicesRepository{Postgres: dbPool})
	groupsService:= service.NewGroupsService(&repository.PostgresGroupsRepository{Postgres: dbPool})

	c := make(chan int)

	defer dbPool.Close()

	devicesService.ManageDevices(mqttClient)
	groupsService.ManageGroups(mqttClient)

	<- c
}
