package model

import (
	"encoding/json"
	"sync/atomic"
)

type UnicornRequestStatus int

const (
	UnicornRequestQueued     UnicornRequestStatus = iota
	UnicornRequestInProgress                      = iota * 10
	UnicornRequestCompleted                       = iota * 10
)

func (u UnicornRequestStatus) String() string {
	switch u {
	case UnicornRequestQueued:
		return "QUEUED"
	case UnicornRequestInProgress:
		return "IN_PROGRESS"
	case UnicornRequestCompleted:
		return "COMPLETED"
	default:
		return "UNICORN_REQUEST_STATUS_UNSPECIFIED"
	}

}

func (u UnicornRequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

type UnicornRequest struct {
	Status          UnicornRequestStatus
	RequestedAmount int
	ReceivedAmount  atomic.Int32
	AvailableAmount atomic.Int32
}

func (u *UnicornRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Status          UnicornRequestStatus `json:"status"`
		RequestedAmount int                  `json:"requested_amount"`
		ReceivedAmount  int                  `json:"received_amount"`
		AvailableAmount int                  `json:"available_amount"`
	}{
		Status:          u.Status,
		RequestedAmount: u.RequestedAmount,
		ReceivedAmount:  int(u.ReceivedAmount.Load()),
		AvailableAmount: int(u.AvailableAmount.Load()),
	})
}

type UnicornRequestId string

func (id UnicornRequestId) String() string {
	return string(id)
}
