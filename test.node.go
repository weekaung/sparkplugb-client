package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"numet.ai/sparkplug/sparkplug"
)

// ******************************************************************************
// ******************************* Node Handlers ********************************
// ******************************************************************************
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("Received message, topic=", msg.Topic())
	p := sparkplug.Payload{}
	err := p.DecodePayload(msg.Payload())
	if err != nil {
		fmt.Println(err)
	}
	ms := p.Metrics
	for i := range ms {
		fmt.Println("Metric: Name=", ms[i].Name, ", DataType=", ms[i].DataType.String(), ", Value=", ms[i].Value)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var reconnectingHandler mqtt.ReconnectHandler = func(client mqtt.Client, options *mqtt.ClientOptions) {
	fmt.Printf("Reconnect handler lost")
	// Note: Need to increment the bdSeq here for reconnecting
	wp, err := sparkplug.GetWillPayload(1)
	if err != nil {
		fmt.Println("Error encoding will payload: ", err)
	}
	options.WillPayload = wp
}

// ******************************************************************************
// ******************************************************************************

func main() {
	// Client Node
	node := sparkplug.ClientNode{
		Config: sparkplug.Config{
			ServerUrl: "192.168.11.61",
			Username:  "da572c34-ca18-11ed-b207-0242ac150002",
			Password:  "yTk92F]G6H7uRn5U>r{D231g#480pjQi",
			ClientID:  "54:36:9b:29:c1:fe", // Device MAC
			GroupID:   "SMART_SECURE",
			NodeID:    "54:36:9b:29:c1:fe",
		},
		MessagePubHandler:   &messagePubHandler,
		ConnectHandler:      &connectHandler,
		ConnectLostHandler:  &connectLostHandler,
		ReconnectingHandler: &reconnectingHandler,
	}

	m1 := sparkplug.Metric{
		Name:     "Node Control/Reboot",
		DataType: sparkplug.TypeBool,
		Value:    "false",
	}
	m2 := sparkplug.Metric{
		Name:     "Manufacturer",
		DataType: sparkplug.TypeString,
		Value:    "Winsonic Electronics",
	}
	m3 := sparkplug.Metric{
		Name:     "io/di/01",
		DataType: sparkplug.TypeInt,
		Value:    "0",
	}
	m4 := sparkplug.Metric{
		Name:     "io/di/02",
		DataType: sparkplug.TypeInt,
		Value:    "0",
	}
	ms := []sparkplug.Metric{}
	ms = append(ms, m1)
	ms = append(ms, m2)
	ms = append(ms, m3)
	ms = append(ms, m4)

	err := node.Connect(0)
	if err != nil {
		fmt.Println(err)
	}
	err = node.SendPayload(sparkplug.MESSAGETYPE_NBIRTH, ms)
	if err != nil {
		fmt.Println(err)
	}
	// Change data and send NData
	m4.Value = "1"
	ms1 := []sparkplug.Metric{}
	ms1 = append(ms1, m4)
	err = node.SendPayload(sparkplug.MESSAGETYPE_NDATA, ms1)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 500)
}
