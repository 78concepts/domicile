package main

import (
	"78concepts.com/domicile/internal/broker"
	"78concepts.com/domicile/internal/database"
	"78concepts.com/domicile/internal/devices"
	"context"
)

func main() {

	// Connect to the MQTT broker
	ctx, ctxCancel:= context.WithCancel(context.Background())

	var client = broker.NewMqttClient(ctx, ctxCancel)

	// Connect to the database
	dbPool:= database.NewPGXPool()

	devicesService:= devices.NewService(&devices.PostgresRepository{Postgres: dbPool})

	// Listen to changes in the available devices
	c := make(chan int)

	defer dbPool.Close()

	devices.HandleDevices(client, devicesService)
	devices.HandleGroups(client, devicesService)

	<- c
}
