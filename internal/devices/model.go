package devices

import "time"

type Device struct {
	IeeeAddress string `json:"ieeeAddr"`
	DateCreated time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateCreated"`
	DateCode string `json:"dateCode"`
	FriendlyName string `json:"friendlyName"`
	AreaId *uint64 `json:"areaId"`
	Description *string `json:"description"`
	Manufacturer string `json:"manufacturerName"`
	Model *string `json:"model"`
	ModelId string `json:"modelID"`
	LastSeen *time.Time `json:"lastSeen"`
	Vendor *string `json:"vendor"`
	Type string `json:"type"`
	Battery *int32 `json:"battery"`
	Active bool `json:"active""`
}

type Area struct {
	Id uint64 `json:"id"`
	DateCreated time.Time `json:"dateCreated"`
	Name string `json:"id"`
}

type TemperatureReport struct {
	DeviceId string `json:"ieeeAddr"`
	AreaId uint64 `json:"areaId"`
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}

type HumidityReport struct {
	DeviceId string `json:"ieeeAddr"`
	AreaId uint64 `json:"areaId"`
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}

type PressureReport struct {
	DeviceId string `json:"ieeeAddr"`
	AreaId uint64 `json:"areaId"`
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}
type IlluminanceReport struct {
	DeviceId string `json:"ieeeAddr"`
	AreaId uint64 `json:"areaId"`
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
	ValueLux float64 `json:"value"`
}