package broker

import (
	"78concepts.com/domicile/internal/config"
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

const (
	TopicRoot = "zigbee2mqtt"
)

type MqttClient struct {
	Ctx        context.Context
	CtxCancel  context.CancelFunc
	Conn mqtt.Client
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v\n", err)
	panic("Connection lost")
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connection to MQTT server successful");
}

var messagePublishHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func NewMqttClient(ctx context.Context, ctxCancel context.CancelFunc) MqttClient {

	configuration := config.GetConfig()

	// Create message broker client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", configuration.Broker.Host, configuration.Broker.Port))
	opts.SetClientID(configuration.Broker.ClientId)
	opts.SetUsername(configuration.Broker.User)
	opts.SetPassword(configuration.Broker.Pass)
	opts.SetDefaultPublishHandler(messagePublishHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	conn := mqtt.NewClient(opts)

	if token := conn.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client:= MqttClient{
		Ctx: ctx,
		CtxCancel: ctxCancel,
		Conn: conn,
	}

	return client
}
