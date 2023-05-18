# Sparkplug B Client for Golang #
## This library complies to the Sparkplug 3.0.0 specification ###
### Note: Only portions of it, just enough for Winsonic DataIO implementation ##

###### To compile the sparkplug proto file
protoc --go_out=. sproto/sparkplug_b.proto

@author David Lee [ Winsonic Electronics ]

# Sample usage in main.go #
## Sparkplug Topic Namespace Elements ##
### namespace/group_id/message_type/edge_node_id/[device_id] ###
### namespace = spBv1.0
### group_id = user defined, logical grouping of Sparkplug edge nodes
### message_type = [NBIRTH, NDEATH, NDATA, NCMD]
### edge_node_id = ID of Sparkplug edge node
### device_id = (optional) ID of attached devices

## Operational Behaviour ##
### 1. Connect to broker using sparkplug.ClientNode
### 2. Send all node information in NBIRTH 
       message.SendPayload(sparkplug.MESSAGETYPE_NBIRTH, ms)
#### Note: sparkplug.ClientNode will automatically send an NDEATH message in the payload header. 
### 3. If node data change, send only changed data using NDATA
       node.SendPayload(sparkplug.MESSAGETYPE_NDATA, ms1)
### 4. Implement the MessagePubHandler to handle NCMD messages
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
### Note: Handle command messages here.
###       If receive Node Control/Rebirth = true, resend NBIRTH message