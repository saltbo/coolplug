package binding

import (
	"mime/multipart"
)

type Plugin struct {
	Name  string                `form:"name" binding:"required"`
	Intro string                `form:"intro" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}
