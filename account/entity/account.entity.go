package entity

type Account struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

type CreateAccountReq struct {
	Owner       string `json:"owner" binding:"required"`
	Currency    string `json:"currency" binding:"required,currency"`
	AccountType string `json:"account_type" binding:"required,oneof=Savings Current"`
}

type GetAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type ListAccountReq struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=2"`
}

type UpdateAccountReq struct {
	Balance int64 `json:"balance"`
}
