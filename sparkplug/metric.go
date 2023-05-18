/*
Sparkplug 3.0.0
Note: Complies to v3.0.0 of the Sparkplug specification
      to the extent needed for Winsonic DataIO and other industrial 4.0 products.
Copyright (c) 2023 Winsonic Electronics, Taiwan
@author David Lee

* This program and the accompanying materials are made available under the
* terms of the Eclipse Public License 2.0 which is available at
* http://www.eclipse.org/legal/epl-2.0.
*/
package sparkplug

import "fmt"

type Metric struct {
	Name     string
	DataType DataType
	// IntValue    int
	// FloatValue  float32
	// BoolValue   bool
	// StringValue string
	Value string
}

type DataType uint32

const (
	TypeInt    DataType = 3
	TypeFloat  DataType = 9
	TypeBool   DataType = 11
	TypeString DataType = 12
)

func (d *DataType) String() string {
	switch *d {
	case TypeInt:
		return "TypeInt"
	case TypeFloat:
		return "TypeFloat"
	case TypeBool:
		return "TypeBool"
	case TypeString:
		return "TypeString"
	}

	fmt.Println(int(d.toUint32()))
	return "error"
}

func (d DataType) toUint32() uint32 {
	return uint32(d)
}
