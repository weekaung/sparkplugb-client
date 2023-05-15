package sparkplug

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ClientNode struct {
	bdSeq               int
	seq                 int
	client              mqtt.Client
	Config              Config
	MessagePubHandler   *mqtt.MessageHandler
	ConnectHandler      *mqtt.OnConnectHandler
	ConnectLostHandler  *mqtt.ConnectionLostHandler
	ReconnectingHandler *mqtt.ReconnectHandler
}

type Config struct {
	ServerUrl string
	Username  string
	Password  string
	ClientID  string
	GroupID   string
	NodeID    string
}

// Connect will connect to the MQTT broker
// Need to provide a bdSeq number which should increment with every connect
// The bdSeq number is from 0 to 255
func (c *ClientNode) Connect(bdSeq int) error {
	// Increment the bdSeq on every connect
	c.bdSeq = bdSeq

	opts := mqtt.NewClientOptions()
	// Set the connection parameters
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.Config.ServerUrl, 1883))
	opts.SetClientID(c.Config.ClientID)
	opts.SetUsername(c.Config.Username)
	opts.SetPassword(c.Config.Password)
	// Set the handlers
	opts.SetDefaultPublishHandler(*c.MessagePubHandler)
	opts.OnConnect = *c.ConnectHandler
	opts.OnConnectionLost = *c.ConnectLostHandler
	opts.OnReconnecting = *c.ReconnectingHandler

	// Set to Auto re-connect
	opts.SetAutoReconnect(true)
	// Set Clean Session on broker to false
	// Broker will remember the previous connection
	opts.CleanSession = false

	// Set the Will topic and message
	opts.WillEnabled = true
	opts.WillQos = 1
	opts.WillRetained = false
	opts.WillTopic = namespace + "/" + c.Config.GroupID + "/" + MESSAGETYPE_NDEATH + "/" + c.Config.NodeID
	// Encode and set Will payload
	wp, err := GetWillPayload(c.bdSeq)
	if err != nil {
		fmt.Println("Error encoding will payload: ", err)
		return err
	}
	opts.WillPayload = wp

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		fmt.Println("Error connecting to MQTT broker: ", err)
		return err
	}

	// Subscribe to receive NCMD messages
	c.subscribeNCMD()
	return nil
}

// GetWillPayload will return the formatted payload with a single metric "bdSeq"
func GetWillPayload(bdSeq int) ([]byte, error) {
	m := Metric{
		Name:     "bdSeq",
		DataType: TypeInt,
		Value:    strconv.Itoa(bdSeq),
	}
	ms := []Metric{}
	ms = append(ms, m)

	p := Payload{
		Metrics: ms,
	}
	return p.EncodePayload(true)
}

func (c *ClientNode) subscribeNCMD() {
	topic := namespace + "/" + c.Config.GroupID + "/" + MESSAGETYPE_NCMD + "/" + c.Config.NodeID
	token := c.client.Subscribe(topic, byte(1), nil)
	token.Wait()
	fmt.Println("Subscribed to topic:", topic)
}

// SendPayload is used to send payload metrics to the MQTT server
// MessageType is either MESSAGETYPE_NBIRTH or MESSAGETYPE_NDATA
// Note: Do not use other message types other than the above 2 types
func (c *ClientNode) SendPayload(messageType string, metrics []Metric) error {
	p := Payload{
		Metrics: metrics,
	}
	// If NBIRTH, add bdSeq & Rebirth metric
	if messageType == MESSAGETYPE_NBIRTH {
		m1 := Metric{
			Name:     "Node Control/Rebirth",
			DataType: TypeBool,
			Value:    "false",
		}
		p.Metrics = append(p.Metrics, m1)

		m2 := Metric{
			Name:     "bdSeq",
			DataType: TypeInt,
			Value:    strconv.Itoa(c.bdSeq),
		}
		p.Metrics = append(p.Metrics, m2)
	}
	// Set the sequence number and increment for next publish
	p.Seq = uint64(c.seq)
	if c.seq == 255 {
		c.seq = 0
	} else {
		c.seq = c.seq + 1
	}
	// Encode payload
	b, err := p.EncodePayload(false)
	if err != nil {
		fmt.Println("Error encoding", messageType, "payload: ", err)
		return err
	}
	topic := namespace + "/" + c.Config.GroupID + "/" + messageType + "/" + c.Config.NodeID
	//fmt.Println("publish:", client, topic, b)
	token := c.client.Publish(topic, 0, false, b)
	token.Wait()
	return nil
}

// Disconnect disconnects the client from the MQTT server
func (c *ClientNode) Disconnect() {
	c.client.Disconnect(0)
}
