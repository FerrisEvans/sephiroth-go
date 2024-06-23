package sys

import "sephiroth-go/model"

type JwtBlacklist struct {
	model.BaseModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
