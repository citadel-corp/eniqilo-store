package user

type StaffResponse struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type CustomerResponse struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}
