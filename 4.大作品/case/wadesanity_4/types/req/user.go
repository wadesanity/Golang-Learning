package types

type UserRegisterReq struct {
	Name   string `form:"name" binding:"required,min=3,max=15" json:"name"`
	Pwd    string `form:"pwd" binding:"required,min=5,max=16" json:"pwd"`
	Avatar string `form:"avatar" binding:"omitempty" json:"avatar"`
}

type UserLoginReq struct {
	Name string `form:"name" binding:"required,min=3,max=15" json:"name"`
	Pwd  string `form:"pwd" binding:"required,min=5,max=16" json:"pwd"`
}

type UserChangePwdReq struct {
	ID     uint   `form:"id" binding:"omitempty" json:"id"`
	PwdOld string `form:"pwd_old" binding:"required,min=5,max=16" json:"pwd_old"`
	PwdNew string `form:"pwd_new" binding:"required,min=5,max=16,nefield=PwdOld" json:"pwd_new"`
}

type UserChangeAvatarReq struct {
	ID     uint   `form:"id" binding:"omitempty" json:"id"`
	Avatar string `form:"avatar" binding:"required,min=5,max=16" json:"avatar"`
}

