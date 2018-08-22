package user

import (
	"cash_flow/util/conn"
	"cash_flow/util/password"
	"fmt"
	"time"
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
	hash, err := password.HashAndSalt(u.Password)
	if err != nil {
		return err
	}
	now := time.Now()
	u.PasswordHash = hash
	u.TimeZone = "UTC"
	err = stmt.QueryRow(
		u.Name,
		u.Email,
		u.PasswordHash,
		nil,
		u.TimeZone,
		now,
	).Scan(&u.Id)
	if err != nil {
		return err
	}
	u.CreatedAt = now
	return nil
}

func (u *User) Update() error {
	var err error
	stmt, err := conn.DB.Prepare(
		`UPDATE users SET
			name = $1,
			email = $2,
			password_hash = $3,
			password_token = $4,
			activation_token = $5,
			activated_at = $6,
			time_zone = $7,
			updated_at = $8
                WHERE deleted_at IS NULL AND id = $9`)
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = stmt.Exec(
		u.Name,
		u.Email,
		u.PasswordHash,
		u.PasswordToken,
		u.ActivationToken,
		u.ActivatedAt,
		u.TimeZone,
		now.Format("2006-01-02 15:04:05"),
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
	now := time.Now()
	_, err = stmt.Exec(
		now.Format("2006-01-02 15:04:05"),
		u.Id,
	)
	if err == nil {
		u.DeletedAt = &now
	}
	return err
}
