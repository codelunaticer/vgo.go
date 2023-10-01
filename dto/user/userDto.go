package userDto

type RegisterDto struct {
	Username string `json:"username" binding:"required,max=13,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginDto struct {
	Username string `json:"username" binding:"required,max=13,min=3"`
	Password string `json:"password" binding:"required"`
}

type SetPassDto struct {
	Username       string `json:"username" binding:"required,max=13,min=3"`
	OldPass        string `json:"old_pass" binding:"required"`
	NewPass        string `json:"new_pass" binding:"required,min=6"`
	ConfirmNewPass string `json:"confirm_new_pass" binding:"required,eqfield=NewPass"`
}

type AddRootAccountDto struct {
	UserId uint `json:"user_id" binding:"required"`
}

type DelRootAccountDto struct {
	UserId uint `json:"user_id" binding:"required"`
}

type DisUserDto struct {
	UserId uint `json:"user_id" binding:"required"`
}

type OpenUserDto struct {
	UserId uint `json:"user_id" binding:"required"`
}

type SetUserInfoDto struct {
	NikeName string `json:"nike_name"` // 别名
	Avatar   string `json:"avatar"`    // 头像
	QQ       string `json:"qq"`        // qq号
	Wechat   string `json:"wechat"`    // 微信号
	Email    string `json:"email"`     // 邮箱号
	Github   string `json:"github"`    // github号
}

// ResetUserPassDto 修改某用户密码
type ResetUserPassDto struct {
	UserId uint `json:"user_id" uri:"user_id" binding:"required"`
}
