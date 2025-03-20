package service

import "unicorn/model"

type UnicornRequestService interface {
	CreateRequest(int) model.UnicornRequestId
	GetRequest(model.UnicornRequestId) (model.UnicornRequest, bool)
}
