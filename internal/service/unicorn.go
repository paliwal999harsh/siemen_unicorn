package service

import (
	"unicorn/pkg/model"
)

type UnicornService interface {
	GetUnicorn(model.UnicornRequestId) []model.Unicorn
}
