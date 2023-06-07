package models

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type User struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	PhoneNumber  int    `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

// GetUserByEmail retrieves a user by their email address from the database.
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT user_id, name, phone_number, email_address, password FROM users WHERE email_address = $1", email).Scan(
		&user.UserID,
		&user.Name,
		&user.PhoneNumber,
		&user.EmailAddress,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Name, &user.PhoneNumber, &user.EmailAddress, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(id string) (User, error) {
	var user User

	err := db.QueryRow("SELECT * FROM users WHERE user_id = $1", id).Scan(&user.UserID, &user.Name, &user.PhoneNumber, &user.EmailAddress, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CreateUser(name string, phoneNumber int, emailAddress, password string) (User, error) {
	var user User

	stmt, err := db.Prepare("INSERT INTO users (name, phone_number, email_address, password) VALUES ($1, $2, $3, $4) RETURNING user_id")
	if err != nil {
		return User{}, err
	}

	err = stmt.QueryRow(name, phoneNumber, emailAddress, password).Scan(&user.UserID)
	if err != nil {
		return User{}, err
	}

	user.Name = name
	user.PhoneNumber = phoneNumber
	user.EmailAddress = emailAddress
	user.Password = password

	return user, nil
}

func UpdateUser(id string, name string, phoneNumber int, emailAddress, password string) error {
	stmt, err := db.Prepare("UPDATE users SET name = $1, phone_number = $2, email_address = $3, password = $4 WHERE user_id = $5")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(name, phoneNumber, emailAddress, password, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("User not found")
	}

	return nil
}

func DeleteUser(id string) error {
	stmt, err := db.Prepare("DELETE FROM users WHERE user_id = $1")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("User not found")
	}

	return nil
}
