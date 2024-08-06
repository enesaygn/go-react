package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sasa-elterminali-service/internal/core/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = "SECRETKEY"

// Login authenticates user and returns JWT token
func (p *DB) Login(req *models.LoginRequest) (*models.LoginResponse, error) {

	var employeeInfo models.EmployeeInfo
	var loginResponse models.LoginResponse

	// req.Password sha256 hash of the password

	hash := sha256.New()
	hash.Write([]byte(req.Password))
	hashSum := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashSum)

	query := `
	SELECT e.employee_id
		,e.sicil_no
		,e.full_name
		,e.e_mail
	FROM techup.employee e
	WHERE (
			e.e_mail = $1
			OR e.sicil_no = $1
			)
		AND e."password" = $2
		AND is_active = true;
`
	err := p.postgres.QueryRow(context.Background(), query, req.UserInfo, hashString).Scan(&employeeInfo.EmployeeID,
		&employeeInfo.SicilNo, &employeeInfo.FullName, &employeeInfo.Email)
	if err != nil {
		return nil, err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["employee_id"] = employeeInfo.EmployeeID
	claims["sicil_no"] = employeeInfo.SicilNo
	claims["full_name"] = employeeInfo.FullName
	claims["e_mail"] = employeeInfo.Email
	claims["exp"] = time.Now().Add(time.Hour * 1000000).Unix()

	t, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}

	loginResponse.Token = t

	return &loginResponse, err

}
