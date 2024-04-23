package model

		import (
			"database/sql"
			"log"

			"time"

			"golang.org/x/crypto/bcrypt"
		)

type User struct {
	ID                      int
	GID                     string
	UserName                string
	Email                   string
	Mobile                  string
	Gender                  string
	DateOfBirth             time.Time
	StartDate               time.Time
	EndDate                 time.Time
	Active                  bool
	UpdatedDate             time.Time
	CreatedDate             time.Time
	LoginUserName           string
	LoginUserPasswordHash   string
	ChangePassword          bool
	MobileVerifiedDate      time.Time
	MobileVerified          bool
	PasswordChangedDate     time.Time
	EmailVerified           bool
	EmailVerifiedDate       time.Time
	EmailVerificationCode   string
	EmailVerificationExpiry time.Time
}

type UserRepository interface {
	AuthenticateUser(username, password string) (*User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (ur *userRepo) AuthenticateUser(username, password string) (*User, error) {
	// Query the database to retrieve the user based on the username
	user := &User{}
	err := ur.db.QueryRow("SELECT id,users_name,login_user_name, login_user_password FROM public.users WHERE login_user_name = $1", username).
		Scan(&user.ID, &user.UserName, &user.LoginUserName, &user.LoginUserPasswordHash)
	if err != nil {
		log.Println("Error retrieving user from database:", err)
		return nil, err
	}

	hashedPassword := []byte(user.LoginUserPasswordHash)
	providedPassword := []byte(password)

	// Compare the hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return nil, nil // Invalid credentials
	}
	return user, nil // Authentication successful
}
