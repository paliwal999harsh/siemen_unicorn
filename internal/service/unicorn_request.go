package service

import "unicorn/pkg/model"

type UnicornRequestService interface {
	CreateRequest(int) model.UnicornRequestId
	GetRequest(model.UnicornRequestId) (*model.UnicornRequest, bool)
}
