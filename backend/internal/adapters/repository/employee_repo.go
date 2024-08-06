package repository

import (
	"context"
	"fmt"
	"sasa-elterminali-service/internal/core/models"

	"github.com/dgrijalva/jwt-go"
)

func (p *DB) GetEmployees() ([]models.Employee, error) {
	query := "SELECT employee_id, full_name, user_name, email FROM techup.employee"
	rows, err := p.postgres.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.EmployeeID, &emp.FullName, &emp.UserName, &emp.Email); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (p *DB) GetEmployee(req *models.RequestBody) (*models.Employee, error) {
	var emp models.Employee
	query := "SELECT employee_id, full_name, user_name, email FROM techup.employee WHERE employee_id=$1"
	err := p.postgres.QueryRow(context.Background(), query, req.EmployeeID).
		Scan(&emp.EmployeeID, &emp.FullName, &emp.UserName, &emp.Email)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (p *DB) CreateEmployee(emp *models.Employee) error {
	query := "INSERT INTO techup.employee (full_name, user_name, email) VALUES ($1, $2, $3)"
	_, err := p.postgres.Exec(context.Background(), query, emp.FullName, emp.UserName, emp.Email)
	return err
}

func (p *DB) UpdateEmployee(emp *models.Employee) error {
	query := "UPDATE techup.employee SET full_name=$1, user_name=$2, email=$3 WHERE employee_id=$4"
	_, err := p.postgres.Exec(context.Background(), query, emp.FullName, emp.UserName, emp.Email, emp.EmployeeID)
	return err
}

func (p *DB) DeleteEmployee(req *models.RequestBody) error {
	query := "DELETE FROM techup.employee WHERE employee_id=$1"
	_, err := p.postgres.Exec(context.Background(), query, req.EmployeeID)
	return err
}

func (p *DB) GetEmployeeIDsByPermission(permissionGroupIDs []int, createUserID int, parameter, status, makineAdi, arizaTipi string,
	pozisyonlar []int) (*models.EmployeesAndCreatedUser, error) {
	var employeesInfo models.EmployeesAndCreatedUser
	var employees []models.Employee
	var createUserName string
	query1 := `
		SELECT er.employee_id
			,full_name::TEXT
		FROM techup.employee_role er
		LEFT JOIN techup.employee e ON e.employee_id = er.employee_id
		WHERE is_deleted = false AND yetki_grup_id = ANY ($1)
	`

	rows, err := p.postgres.Query(context.Background(), query1, permissionGroupIDs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var employee models.Employee
		if err := rows.Scan(&employee.EmployeeID, &employee.FullName); err != nil {
			return nil, err
		}

		employees = append(employees, employee)

	}

	query2 := `
		SELECT full_name
		FROM techup.employee e
		WHERE e.employee_id = $1
	`

	err = p.postgres.QueryRow(context.Background(), query2, createUserID).Scan(&createUserName)
	if err != nil {
		return nil, err
	}

	if status == ARIZA_YENI || status == ARIZA_ISLEMDE || status == ARIZA_PROSES_ONAYI || status == ARIZA_TEST_ONAYI ||
		status == ARIZA_TEKRAR || status == ARIZA_OLUMSUZ_TEST || status == ARIZA_TEST_TALEBI {

		employeesInfo.Parameter = parameter
		employeesInfo.Message = fmt.Sprintf("Makine %s / %v pozisyon/ları \"%s\"", makineAdi, pozisyonlar, arizaTipi)

	} else if ARIZA_YAPILDI == status {

		employeesInfo.Parameter = parameter
		employeesInfo.Message = fmt.Sprintf("Makine %s / %v pozisyon/ları", makineAdi, pozisyonlar)
	}

	employeesInfo.Employees = employees
	employeesInfo.FullName = createUserName

	return &employeesInfo, nil
}

func (p *DB) GetUserRoleFromToken(req *models.GetUserRoleFromTokenRequest) (*models.GetRoleRGetUserRoleFromTokenResponse, error) {
	var role models.GetRoleRGetUserRoleFromTokenResponse
	var employeeIDSTR string

	// Parse and validate the token
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the key for validation
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token claims are not of type jwt.MapClaims")
	}

	employeeIDSTR = fmt.Sprintf("%v", claims["employee_id"])
	print(employeeIDSTR)

	query := `
		SELECT employee_id, yetki_grubu_adi
		FROM techup.employee_role er
		LEFT JOIN techup.permission_group p ON er.yetki_grup_id = p.id
		WHERE employee_id = $1::int AND er.yetki_grup_id in (116,117,118,119,120)
	`
	err = p.postgres.QueryRow(context.Background(), query, req.Token).Scan(&role.EmployeeId, &role.Role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
