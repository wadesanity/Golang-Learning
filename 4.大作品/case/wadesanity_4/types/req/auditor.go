package types

type AuditorReq struct {
	ID     uint  `json:"id" form:"id" binding:"required"`
	Status *uint `json:"status" form:"status" binding:"required"`
}
