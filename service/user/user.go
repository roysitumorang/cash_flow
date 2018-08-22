package user

import (
	"cash_flow/util/conn"
	"cash_flow/util/crypt"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

func New() *User {
	return new(User)
}

func Find(userId int) (*User, error) {
	var u = New()
	err := conn.DB.QueryRow(
		`SELECT
			id,
			name,
			email,
			password_hash,
			password_token,
			activation_token,
			activated_at,
			time_zone,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL AND id = $1`,
		userId,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PasswordToken,
		&u.ActivationToken,
		&u.ActivatedAt,
		&u.TimeZone,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

func FindByActivationToken(token string) (*User, error) {
	var u = New()
	err := conn.DB.QueryRow(
		`SELECT
			id,
			name,
			email,
			password_hash,
			password_token,
			activation_token,
			activated_at,
			time_zone,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL AND activation_token = $1`,
		token,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PasswordToken,
		&u.ActivationToken,
		&u.ActivatedAt,
		&u.TimeZone,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

func FindByPasswordToken(token string) (*User, error) {
	var u = New()
	err := conn.DB.QueryRow(
		`SELECT
			id,
			name,
			email,
			password_hash,
			password_token,
			activation_token,
			activated_at,
			time_zone,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL AND password_token = $1`,
		token,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PasswordToken,
		&u.ActivationToken,
		&u.ActivatedAt,
		&u.TimeZone,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

func FindByEmail(email string) (*User, error) {
	var u = New()
	err := conn.DB.QueryRow(
		`SELECT
			id,
			name,
			email,
			password_hash,
			password_token,
			activation_token,
			activated_at,
			time_zone,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL AND email = $1`,
		email,
	).Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.PasswordToken,
		&u.ActivationToken,
		&u.ActivatedAt,
		&u.TimeZone,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

func Authenticate(email, password string) (string, error) {
	var (
		err error
	)
	u, err := FindByEmail(email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["UserId"] = u.Id
	claims["exp"] = time.Now().UTC().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, err
}

func Where(term string, page int) (Users, error) {
	var err error
	users := Users{}
	if page == 0 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit
	query := fmt.Sprintf("%%%s%%", term)
	rows, err := conn.DB.Query(
		`SELECT
			id,
			name,
			email,
			password_hash,
			password_token,
			activation_token,
			activated_at,
			time_zone,
			created_at,
			updated_at,
			deleted_at
		FROM users
		WHERE deleted_at IS NULL AND (name ILIKE $1 OR email ILIKE $2)
		LIMIT $3 OFFSET $4`,
		query,
		query,
		limit,
		offset,
	)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		u := New()
		err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.PasswordHash,
			&u.PasswordToken,
			&u.ActivationToken,
			&u.ActivatedAt,
			&u.TimeZone,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	return users, err
}

func (u *User) Validate() bool {
	u.Errors = make(map[string]string)
	re := regexp.MustCompile(".+@.+\\..+")
	if strings.TrimSpace(u.Name) == "" {
		u.Errors["name"] = "blank name"
	}
	if strings.TrimSpace(u.Email) == "" {
		u.Errors["email"] = "blank email address"
	} else if matched := re.Match([]byte(u.Email)); !matched {
		u.Errors["email"] = "invalid email address"
	} else {
		var conflict int
		if u.Id == nil {
			stmt, _ := conn.DB.Prepare(`SELECT COUNT(1) FROM users WHERE email = $1`)
			stmt.QueryRow(u.Email).Scan(&conflict)
		} else {
			stmt, _ := conn.DB.Prepare(`SELECT COUNT(1) FROM users WHERE email = $1 AND id != $2`)
			stmt.QueryRow(u.Email, u.Id).Scan(&conflict)
		}
		if conflict > 0 {
			u.Errors["email"] = "email exists"
		}
	}
	if u.Id == nil {
		if strings.TrimSpace(u.Password) == "" {
			u.Errors["password"] = "blank password"
		} else if utf8.RuneCountInString(strings.TrimSpace(u.PasswordConfirmation)) < 6 {
			u.Errors["password"] = "minimum 6 characters"
		} else if strings.TrimSpace(u.PasswordConfirmation) == "" {
			u.Errors["password"] = "blank password confirmation"
		} else if u.Password != u.PasswordConfirmation {
			u.Errors["password"] = "invalid password confirmation"
		}
	}
	return len(u.Errors) == 0
}

func (u *User) Create() error {
	stmt, err := conn.DB.Prepare(
		`INSERT INTO users (
			name,
			email,
			password_hash,
			activation_token,
			time_zone,
			created_at
                ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	activationToken, err := crypt.GenerateRandomString(64)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	u.PasswordHash = string(hash)
	u.ActivationToken = &activationToken
	u.TimeZone = "UTC"
	err = stmt.QueryRow(
		u.Name,
		u.Email,
		u.PasswordHash,
		u.ActivationToken,
		u.TimeZone,
		now,
	).Scan(&u.Id)
	if err != nil {
		return err
	}
	u.CreatedAt = now
	return nil
}

func (u *User) Activate() error {
	var err error
	stmt, err := conn.DB.Prepare(
		`UPDATE users SET
			activation_token = $1,
			activated_at = $2,
			updated_at = $3
                WHERE deleted_at IS NULL AND activation_token IS NOT NULL AND activated_at IS NULL AND id = $4`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = stmt.Exec(
		nil,
		now,
		now,
		u.Id,
	)
	if err == nil {
		u.ActivationToken = nil
		u.ActivatedAt = &now
		u.UpdatedAt = &now
	}
	return err
}

func (u *User) ResetPassword() error {
	var err error
	stmt, err := conn.DB.Prepare(
		`UPDATE users SET
			password_token = $1,
			updated_at = $2
                WHERE deleted_at IS NULL AND id = $3`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	passwordToken, err := crypt.GenerateRandomString(64)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		passwordToken,
		now,
		u.Id,
	)
	if err == nil {
		u.PasswordToken = &passwordToken
		u.UpdatedAt = &now
	}
	return err
}

func (u *User) SavePassword() error {
	var err error
	stmt, err := conn.DB.Prepare(
		`UPDATE users SET
			password_token = $1,
			password_hash = $2,
			updated_at = $3
                WHERE deleted_at IS NULL AND id = $4`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		nil,
		hash,
		now,
		u.Id,
	)
	if err == nil {
		u.PasswordToken = nil
		u.PasswordHash = string(hash)
		u.UpdatedAt = &now
	}
	return err
}

func (u *User) Update() error {
	var err error
	stmt, err := conn.DB.Prepare(
		`UPDATE users SET
			name = $1,
			email = $2,
			password_hash = $3,
			time_zone = $4,
			updated_at = $5
                WHERE deleted_at IS NULL AND id = $6`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hash)
	}
	_, err = stmt.Exec(
		u.Name,
		u.Email,
		u.PasswordHash,
		u.TimeZone,
		now,
		u.Id,
	)
	if err == nil {
		u.UpdatedAt = &now
	}
	return err
}

func (u *User) Destroy() error {
	var err error
	stmt, err := conn.DB.Prepare(`UPDATE users SET deleted_at = $1 WHERE deleted_at IS NULL AND id = $2`)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = stmt.Exec(
		now,
		u.Id,
	)
	if err == nil {
		u.DeletedAt = &now
	}
	return err
}
