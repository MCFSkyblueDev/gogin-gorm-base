package dtos

type UpdateUserInfoRequest struct {
	Name *string `form:"name" binding:"required,min=2,max=20"`
	Bio  *string `form:"bio" binding:"omitempty,min=12"`
	// Avatar       *string `json:"avatar"`
	MemberNumber *string `form:"member_number" binding:"omitempty"`
}

type UserInfoResponse struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Bio          *string `json:"bio"`
	Avatar       *string `json:"avatar"`
	MemberNumber *string `json:"member_number"`
}
