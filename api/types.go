package api

import (
	"encoding/json"
	"fmt"
	"net"
)

// Device is the JSON representation of a device from the Kentik API device endpoints.
type Device struct {
	ID          int      `json:"id,string"`
	Name        string   `json:"device_name"`
	Type        string   `json:"device_type"`
	Subtype     string   `json:"device_subtype"`
	Description string   `json:"device_description"`
	IP          net.IP   `json:"ip"`
	SampleRate  int      `json:"device_sample_rate,string"`
	BgpType     string   `json:"device_bgp_type"`
	Plan        Plan     `json:"plan"`
	CdnAttr     string   `json:"cdn_attr"`
	MaxFlowRate int      `json:"max_flow_rate"`
	CompanyID   int      `json:"company_id,string"`
	Customs     []Column `json:"custom_column_data,omitempty"`
}

// Plan is the JSON representation of a Kentik License Plan that describes the services provided by Kentik (i.e.
// permitted data, maximum number of devices, maximum number of FPS, etc.).
type Plan struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// Column is the JSON representation of a custom column from the Kentik API device endpoints.
type Column struct {
	ID   uint64 `json:"field_id,string"`
	Name string `json:"col_name"`
	Type string `json:"col_type"`
}

// SiteAndDeviceCreate is the payload JSON representation wrapper for creating a Site and Device.
type SiteAndDeviceCreate struct {
	Site   *SiteCreate   `json:"site"`
	Device *DeviceCreate `json:"device"`
}

// SiteCreate is the payload JSON representation wrapper for creating a Site.
type SiteCreate struct {
	Title   string `json:"title"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
}

// DeviceCreate is the payload JSON representation wrapper for creating a Device.
type DeviceCreate struct {
	Name        string `json:"device_name"`
	Type        string `json:"device_type"`
	Subtype     string `json:"device_subtype"`
	Description string `json:"device_description"`
	SampleRate  int    `json:"device_sample_rate,string"`
	BgpType     string `json:"device_bgp_type"`
	PlanID      int    `json:"plan_id,omitempty"`
	SiteID      int    `json:"site_id,omitempty"`

	CdnAttr string `json:"cdn_attr"`

	// IPs is the associated sending IP Address(es) for the device. For devices that do not have a set of permanent IP
	// addresses (i.e. cloud devices), this should not be set and AllowNoIP should be set to true.
	IPs []net.IP `json:"sending_ips"`

	// AllowNoIP is a flag to bypass the requirement of at least one IP address for creation of a device. This is only
	// applicable for devices that do not have a set of permanent IP addresses, such as cloud devices.
	AllowNoIP bool `json:"-"`

	// ExportId the optionally associated Cloud Exporter which first received flow and is creating the device.
	ExportId int `json:"cloud_export_id,omitempty"`

	// Region is the optionally associated cloud region for this device
	Region string `json:"cloud_region,omitempty"`

	// Zone is the optionally associated cloud zone for this device
	Zone string `json:"cloud_zone,omitempty"`
}

// DeviceWrapper is a wrapper around the Device, for endpoints that return a single Device nested away next to other
// information.
type DeviceWrapper struct {
	Device *Device `json:"device"`
}

// AllDeviceWrapper is a wrapper around a list of Device instances, for endpoints that return a one or more Devices
// nested in a JSON structure.
type AllDeviceWrapper struct {
	Devices []*Device `json:"devices"`
}

// DevicesFilter is a set of arguments that can be used to filter the devices queried along with the hydration of the
// resulting JSON.
type DevicesFilter struct {

	// FilterCloud restricts the results to devices that were NOT created by cloud exporters.
	FilterCloud bool `schema:"filterCloud"`

	// CloudOnly restricts the results to ONLY devices created by cloud exporters.
	CloudOnly bool `schema:"cloudOnly"`

	// Subtypes restricts the results to only specific device subtypes. An empty or nil instance indicates no filtering
	// of results.
	Subtypes []string `schema:"subtypes"`

	// AugmentWith specifies which additional, potentially large, payload fields should be hydrated. Some supported
	// options are:
	//   - customColumns -
	//   - plan - The associated Plan for the device
	AugmentWith []string `schema:"augmentWith"`

	// Columns specifies additional, metadata fields i.e. table columns, to retrieve as part of the query. Examples
	// include id, device_name, etc.
	Columns []string `schema:"columns"`
}

type Interface struct {
	ID      uint64 `json:"id,string"`
	Index   uint64 `json:"snmp_id,string"`
	Alias   string `json:"snmp_alias"`
	Desc    string `json:"interface_description"`
	Address string `json:"interface_ip"`
	Netmask string `json:"interface_ip_netmask"`
	Addrs   []Addr `json:"secondary_ips"`
}

type InterfaceUpdate struct {
	Index   uint64 `json:"index,string"`
	Alias   string `json:"alias"`
	Desc    string `json:"desc"`
	Speed   uint64 `json:"speed"`
	Type    uint64 `json:"type"`
	Address string `json:"address"`
	Netmask string `json:"netmask"`
	Addrs   []Addr `json:"alias_address"`
}

type Addr struct {
	Address string `json:"address"`
	Netmask string `json:"netmask"`
}

// ClientID creates the encoding of Client ID which may be used for authentication purposes for sending data on behalf
// of the device.
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

func (c *DeviceCreate) NormalizeName() {
	c.Name = NormalizeName(c.Name)
}
