package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DeviceResponse struct {
	Device Device `json:"device"`
}

type Device struct {
	ID          int           `json:"id,string"`
	Name        string        `json:"device_name"`
	MaxFlowRate int           `json:"max_flow_rate"`
	CompanyID   int           `json:"company_id,string"`
	Customs     CustomColumns `json:"custom_columns"`
}

type CustomColumns map[string]uint64

func (d *Device) ClientID() string {
	return fmt.Sprintf("%d:%s:%d", d.CompanyID, d.Name, d.ID)
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

func (c *CustomColumns) MarshalJSON() ([]byte, error) {
	kvs := make([]string, 0, len(*c))
	for k, v := range *c {
		kvs = append(kvs, fmt.Sprintf("%s=%d", k, v))
	}
	return []byte(`"` + strings.Join(kvs, ",") + `"`), nil
}

var split = regexp.MustCompile(`([\w-]+)=(\d+),?`)
