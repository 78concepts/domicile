package model

import "time"

type Device struct {
	IeeeAddress string `json:"ieeeAddr"`
	DateCreated time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
	DateCode *string `json:"dateCode"`
	FriendlyName string `json:"friendlyName"`
	AreaId *uint64 `json:"areaId"`
	Description *string `json:"description"`
	Manufacturer *string `json:"manufacturerName"`
	Model *string `json:"model"`
	ModelId *string `json:"modelID"`
	LastSeen *time.Time `json:"lastSeen"`
	Vendor *string `json:"vendor"`
	Type *string `json:"type"`
	Battery *int32 `json:"battery"`
	Active bool `json:"active""`
}
