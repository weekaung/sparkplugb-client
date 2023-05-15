package main

import (
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/weekaung/sparkplugb-client/sparkplug"
)

// ******************************************************************************
// *************************** Application Handlers *****************************
// ******************************************************************************
var messagePubHandlerA mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Application: Received message by, topic=", msg.Topic())
	topic := strings.Split(msg.Topic(), "/")
	if len(topic) >= 2 && topic[0] == "spBv1.0" {
		if topic[1] == "STATE" {
			fmt.Println("Application: Payload=", string(msg.Payload()))
		} else {
			p := sparkplug.Payload{}
			err := p.DecodePayload(msg.Payload())
			if err != nil {
				fmt.Println(err)
			}
			ms := p.Metrics
			for i := range ms {
				fmt.Println("Application: Metric: Name=", ms[i].Name, ", DataType=", ms[i].DataType.String(), ", Value=", ms[i].Value)
			}
		}
	}
}

var connectHandlerA mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Application: Connected")
}

var connectLostHandlerA mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Application: Connect lost: %v", err)
}

// ******************************************************************************
// ******************************************************************************

func main() {
	// App Node
	app := sparkplug.ClientApp{
		Auth: sparkplug.Auth{
			ServerUrl: "192.168.11.61",
			Username:  "DMS",
			Password:  "12345678901234567890123456789012",
			GroupID:   "DMS",
		},
		MessagePubHandler:  &messagePubHandlerA,
		ConnectHandler:     &connectHandlerA,
		ConnectLostHandler: &connectLostHandlerA,
		// ReconnectingHandler: &reconnectingHandlerA,
	}
	app.Connect()
	app.SetOnline()

	time.Sleep(time.Second * 500)
}
