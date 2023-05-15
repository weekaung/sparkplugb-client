package sparkplug

import (
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
	"numet.ai/sparkplug/sproto"
)

const namespace = "spBv1.0"
const state = "STATE"
const MESSAGETYPE_NBIRTH = "NBIRTH"
const MESSAGETYPE_NDEATH = "NDEATH"
const MESSAGETYPE_NDATA = "NDATA"
const MESSAGETYPE_NCMD = "NCMD"

type Payload struct {
	Timestamp time.Time
	Seq       uint64
	Metrics   []Metric
}

func (p *Payload) EncodePayload(isDeathPayload bool) ([]byte, error) {
	now := time.Now().UnixMilli()
	ms := []*sproto.Payload_Metric{}

	for i, m := range p.Metrics {
		sm := sproto.Payload_Metric{}
		sm.Name = &p.Metrics[i].Name
		dt := m.DataType.toUint32()
		sm.Datatype = &dt
		/**********************************
		TypeInt DataType = 10
		TypeFloat DataType = 12
		TypeBool DataType = 14
		TypeString DataType = 15
		**********************************/
		switch m.DataType {
		case TypeInt:
			iv, err := strconv.ParseUint(m.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(iv)}
		case TypeFloat:
			fv, err := strconv.ParseFloat(m.Value, 32)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_FloatValue{FloatValue: float32(fv)}
		case TypeBool:
			bv, err := strconv.ParseBool(m.Value)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_BooleanValue{BooleanValue: bv}
		case TypeString:
			sm.Value = &sproto.Payload_Metric_StringValue{StringValue: m.Value}
		}
		//fmt.Println(*sm.Name, "=", sm.Value)
		ms = append(ms, &sm)
	}
	//fmt.Println("---------")

	sp := sproto.Payload{}

	if !isDeathPayload {
		// Set Payload timestamp
		tn := uint64(now)
		sp.Timestamp = &tn
		// Set Payload sequence
		sp.Seq = &p.Seq
	}
	sp.Metrics = ms
	return proto.Marshal(&sp)
}

func (p *Payload) DecodePayload(bytes []byte) error {
	pl := sproto.Payload{}
	proto.Unmarshal(bytes, &pl)
	fmt.Println("Payload is ", pl.String())
	if pl.Timestamp != nil {
		p.Timestamp = time.UnixMilli(int64(*pl.Timestamp))
	}
	p.Metrics = make([]Metric, len(pl.Metrics))
	for i := range pl.Metrics {
		p.Metrics[i].Name = *pl.Metrics[i].Name
		p.Metrics[i].DataType = DataType(*pl.Metrics[i].Datatype)
		// Set the Value according to DataType
		switch p.Metrics[i].DataType {
		case TypeInt:
			p.Metrics[i].Value = strconv.FormatUint(uint64(pl.Metrics[i].GetIntValue()), 10)
		case TypeFloat:
			p.Metrics[i].Value = fmt.Sprintf("%f", pl.Metrics[i].GetFloatValue())
		case TypeBool:
			p.Metrics[i].Value = strconv.FormatBool(pl.Metrics[i].GetBooleanValue())
		case TypeString:
			p.Metrics[i].Value = pl.Metrics[i].GetStringValue()
		}
	}
	return nil
}
