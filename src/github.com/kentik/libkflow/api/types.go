package api

import (
	"fmt"
	"regexp"
	"strconv"
)

type DeviceResponse struct {
	Device Device `json:"device"`
}

type Device struct {
	ID        string        `json:"id"`
	Name      string        `json:"device_name"`
	CompanyID string        `json:"company_id"`
	Custom    CustomColumns `json:"custom_columns"`
}

type CustomColumns map[string]uint64

func (d *Device) ClientID() string {
	return fmt.Sprintf("%s:%s:%s", d.CompanyID, d.Name, d.ID)
}

func (c *CustomColumns) UnmarshalJSON(data []byte) error {
	m := map[string]uint64{}
	for _, match := range split.FindAllSubmatchIndex(data, -1) {
		key := string(data[match[2]:match[3]])
		val := string(data[match[4]:match[5]])
		m[key], _ = strconv.ParseUint(val, 10, 64)
	}
	*c = m
	return nil
}

var split = regexp.MustCompile(`(\w+)=(\d+),?`)
