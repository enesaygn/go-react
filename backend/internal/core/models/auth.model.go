package models

type LoginRequest struct {
	UserInfo string `json:"userInfo"`
	Password string `json:"password"`
}

type EmployeeInfo struct {
	EmployeeID int    `json:"employeeID"`
	SicilNo    string `json:"sicilNo"`
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
