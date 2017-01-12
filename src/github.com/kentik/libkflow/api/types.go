package api

import (
	"encoding/json"
	"fmt"
)

type DeviceResponse struct {
	Device Device `json:"device"`
}

type Device struct {
	ID          int      `json:"id,string"`
	Name        string   `json:"device_name"`
	MaxFlowRate int      `json:"max_flow_rate"`
	CompanyID   int      `json:"company_id,string"`
	Customs     []Column `json:"custom_column_data,omitempty"`
}

type Column struct {
	ID   uint64 `json:"field_id,string"`
	Name string `json:"col_name"`
	Type string `json:"col_type"`
}

func (d *Device) ClientID() string {
	return fmt.Sprintf("%d:%s:%d", d.CompanyID, d.Name, d.ID)
}

func (c *Column) UnmarshalFlag(value string) error {
	return json.Unmarshal([]byte(value), c)
}

func (c Column) MarshalFlag() (string, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
