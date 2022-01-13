package model

import "time"

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

