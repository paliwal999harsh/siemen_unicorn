package service

import (
	"unicorn/model"
)

type UnicornService interface {
	GetUnicorn(model.UnicornRequestId) []model.Unicorn
}
