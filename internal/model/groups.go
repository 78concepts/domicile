package model

import "time"

type Group struct {
	Id uint64 `json:"id"`
	DateCreated time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
	FriendlyName string `json:"friendlyName"`
	Active bool `json:"active""`
}

type GroupMember struct {
	GroupId uint64 `json:"groupId"`
	IeeeAddress string `json:"ieeeAddress"`
}

