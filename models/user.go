package models

import (
	"golang-api/services"
	"golang-api/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

// User implements Model
type User struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	Firstname       string     `json:"firstname" gorm:"type:varchar(100)"`
	Lastname        string     `json:"lastname" gorm:"type:varchar(100)"`
	Username        string     `json:"username" gorm:"type:varchar(100)"`
	Email           string     `json:"email" gorm:"type:varchar(100)"`
	Password        string     `json:"password"`
	ProfileImageUrl string     `json:"profile_image_url" gorm:"type:varchar(255)"`
	Roles           []string   `json:"roles" gorm:"json"`
	Verified        bool       `json:"verified" gorm:"default:false"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type CreateUserDto struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUserDto struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type SanitizedUser struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	Email           string    `json:"email"`
	Roles           []string  `json:"roles"`
	Verified        bool      `json:"verified"`
	ProfileImageUrl string    `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LoginUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserInfo struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	ProfileImageUrl string `json:"profile_image_url"`
}

func CountStatsUsers(ws *services.Websocket) map[string]int {
	totalUsers, _ := CountUsers()

	loggedClients := ws.FilterClient(func(c *services.Client) bool {
		return c.Get("logged").(bool)
	})

	annonClients := ws.FilterClient(func(c *services.Client) bool {
		return !c.Get("logged").(bool)
	})

	return map[string]int{
		"loggedUsers": len(loggedClients),
		"annonUsers":  len(annonClients),
		"totalUsers":  int(totalUsers),
	}
}

func FindAllUsers(query services.QueryFilter) ([]User, error) {
	var users []User

	value := DB.Model(&User{}).
		Where("deleted_at IS NULL").
		Offset(query.GetSkip()).
		Where(query.GetWhere()).
		Order(query.GetSort())

	if query.GetLimit() != 0 {
		value.Limit(query.GetLimit())
	}

	err := value.Find(&users).Error

	return users, err
}

func CountUsersByEmail(email string) (int64, error) {
	var count int64

	err := DB.Model(&User{}).
		Where("deleted_at IS NULL").
		Where("email = ?", email).
		Count(&count).Error

	return count, err
}

func CountUsersByUsername(username string) (int64, error) {
	var count int64

	err := DB.Model(&User{}).
		Where("deleted_at IS NULL").
		Where("username", username).
		Count(&count).Error

	return count, err
}

func CountUsers() (int64, error) {
	var count int64

	err := DB.Model(&User{}).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

func (u *User) Sanitize() SanitizedUser {
	return SanitizedUser{
		ID:              u.ID,
		Username:        u.Username,
		Firstname:       u.Firstname,
		Lastname:        u.Lastname,
		Email:           u.Email,
		Roles:           u.Roles,
		Verified:        u.Verified,
		ProfileImageUrl: u.ProfileImageUrl,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func (u *User) IsRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) FindOneById(id int) error {
	return DB.Model(&User{}).
		First(&u, id).Error
}

func (u *User) FindOne(key string, value any) error {
	return DB.Model(&User{}).
		Where("deleted_at IS NULL").
		Where(key, value).
		First(&u).Error
}

func (u *User) Delete() error {
	if err := DeleteTokensByUserID(u.ID); err != nil {
		return err
	}

	del := time.Now()
	u.DeletedAt = &del
	u.Username = "[deleted]"
	u.Email = "[deleted]"
	u.Firstname = "[deleted]"
	u.Lastname = "[deleted]"
	u.ProfileImageUrl = "https://api.dicebear.com/9.x/initials/svg?seed=deleted"
	if err := u.HashPassword(utils.GenerateString(20)); err != nil {
		return err
	}

	return u.Save()
}

func (u *User) Save() error {
	return DB.Save(&u).Error
}

func ClearUsers() error {
	return DB.Exec("DELETE FROM users").Error
}

func UserExists(id int) bool {
	var count int64
	DB.Model(&User{}).Where("id", id).Count(&count)
	return count > 0
}
