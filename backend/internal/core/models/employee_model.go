package models

type EmployeesAndCreatedUser struct {
	FullName  string
	Message   string
	Parameter string
	Employees []Employee
}

type Employee struct {
	EmployeeID int
	FullName   string
	UserName   string
	Email      *string
}

type RequestBody struct {
	EmployeeID int
}

type GetUserRoleFromTokenRequest struct {
	Token string
}

type GetRoleRGetUserRoleFromTokenResponse struct {
	Role       string
	EmployeeId int
}
