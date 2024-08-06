package models

import (
	"golang-api/services"
	"strconv"
	"testing"

	faker "github.com/jaswdr/faker/v2"
	"golang.org/x/crypto/bcrypt"
)

var fake = faker.New()

func TestUserSave(t *testing.T) {
	var nbUsers int = 10
	var result []User
	var users []User

	if err := ConnectDB(true); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := CloseDB(); err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	for i := 0; i < nbUsers; i++ {
		user := User{
			Firstname: fake.Person().FirstName(),
			Lastname:  fake.Person().LastName(),
			Username:  fake.Person().Name(),
			Email:     fake.Internet().Email(),
			Password:  fake.Internet().Password(),
			Roles:     []string{"user"},
		}
		users = append(users, user)
		if err := user.Save(); err != nil {
			t.Error(err)
			return
		}
	}

	if err := DB.Find(&result).Error; err != nil {
		t.Error("Failed to find users:", err)
		return
	}

	if len(users) != nbUsers {
		t.Error("Expected 10 users, got", len(users))
	}

	if len(result) != nbUsers {
		t.Error("Expected 10 users, got", len(result))
	}

	for i, res := range result {
		if res.ID == 0 {
			t.Error("Expected user ID to be set")
		}

		if res.Firstname != users[i].Firstname {
			t.Errorf("Expected %s, got %s", users[i].Firstname, res.Firstname)
		}

		if res.Lastname != users[i].Lastname {
			t.Errorf("Expected %s, got %s", users[i].Lastname, res.Lastname)
		}

		if res.Username != users[i].Username {
			t.Errorf("Expected %s, got %s", users[i].Username, res.Username)
		}

		if res.Email != users[i].Email {
			t.Errorf("Expected %s, got %s", users[i].Email, res.Email)
		}

		if res.Password != users[i].Password {
			t.Errorf("Expected %s, got %s", users[i].Password, res.Password)
		}

		if len(res.Roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(res.Roles))
		}

		if res.Roles[0] != "user" {
			t.Errorf("Expected role 'user', got %s", res.Roles[0])
		}
	}
}

func TestFindAllUsers(t *testing.T) {
	var nbUsers int = 10
	var users []User

	if err := ConnectDB(true); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := CloseDB(); err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	for i := 0; i < nbUsers; i++ {
		user := User{
			Firstname: fake.Person().FirstName(),
			Lastname:  fake.Person().LastName(),
			Username:  fake.Person().Name(),
			Email:     fake.Internet().Email(),
			Password:  fake.Internet().Password(),
			Roles:     []string{"user"},
		}
		users = append(users, user)
		if err := DB.Create(&user).Error; err != nil {
			t.Error("Failed to create user:", err)
			return
		}
	}

	params := map[string][]string{
		"page":  {"1"},
		"limit": {strconv.Itoa(nbUsers)},
	}

	query, err := services.NewQueryFilter(params)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := FindAllUsers(query)
	if err != nil {
		t.Error(err)
		return
	}

	if len(users) != nbUsers {
		t.Errorf("Expected 10 users, got %d", len(users))
	}

	if len(result) != nbUsers {
		t.Errorf("Expected 10 users, got %d", len(result))
	}

	for i, res := range result {
		if res.ID == 0 {
			t.Error("Expected user ID to be set")
		}

		if res.Firstname != users[i].Firstname {
			t.Errorf("Expected %s, got %s", users[i].Firstname, res.Firstname)
		}

		if res.Lastname != users[i].Lastname {
			t.Errorf("Expected %s, got %s", users[i].Lastname, res.Lastname)
		}

		if res.Username != users[i].Username {
			t.Errorf("Expected %s, got %s", users[i].Username, res.Username)
		}

		if res.Email != users[i].Email {
			t.Errorf("Expected %s, got %s", users[i].Email, res.Email)
		}

		if res.Password != users[i].Password {
			t.Errorf("Expected %s, got %s", users[i].Password, res.Password)
		}

		if len(res.Roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(res.Roles))
		}

		if res.Roles[0] != "user" {
			t.Errorf("Expected role 'user', got %s", res.Roles[0])
		}
	}
}

func TestUserFindOneById(t *testing.T) {
	var user User
	var result User

	if err := ConnectDB(true); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := CloseDB(); err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	user = User{
		Firstname: fake.Person().FirstName(),
		Lastname:  fake.Person().LastName(),
		Username:  fake.Person().Name(),
		Email:     fake.Internet().Email(),
		Password:  fake.Internet().Password(),
		Roles:     []string{"user"},
	}

	if err := DB.Create(&user).Error; err != nil { // Check the error on DB.Create
		t.Error("Failed to create user:", err)
		return
	}

	if err := result.FindOneById(user.ID); err != nil {
		t.Error(err)
		return
	}

	if result.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if result.Firstname != user.Firstname {
		t.Errorf("Expected %s, got %s", user.Firstname, result.Firstname)
	}

	if result.Lastname != user.Lastname {
		t.Errorf("Expected %s, got %s", user.Lastname, result.Lastname)
	}

	if result.Username != user.Username {
		t.Errorf("Expected %s, got %s", user.Username, result.Username)
	}

	if result.Email != user.Email {
		t.Errorf("Expected %s, got %s", user.Email, result.Email)
	}

	if result.Password != user.Password {
		t.Errorf("Expected %s, got %s", user.Password, result.Password)
	}

	if len(result.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(result.Roles))
	}

	if result.Roles[0] != "user" {
		t.Errorf("Expected role 'user', got %s", result.Roles[0])
	}
}

func TestCountUsersByEmail(t *testing.T) {
	var user User
	var count int64
	var nbUsers int64 = 5
	var i int64
	var testEmail string = fake.Internet().Email()

	if err := ConnectDB(true); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := CloseDB(); err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	for i = 0; i < nbUsers; i++ {
		user = User{
			Firstname: fake.Person().FirstName(),
			Lastname:  fake.Person().LastName(),
			Username:  fake.Person().Name(),
			Email:     fake.Internet().Email(),
			Password:  fake.Internet().Password(),
			Roles:     []string{"user"},
		}

		if i == 0 {
			user.Email = testEmail
		}

		DB.Create(&user)
	}

	count, err := CountUsersByEmail(testEmail)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestCountUsersByUsername(t *testing.T) {
	var user User
	var count int64
	var nbUsers int64 = 5
	var i int64
	var testUsername string = fake.Person().Name()

	if err := ConnectDB(true); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := CloseDB(); err != nil {
			t.Error("Failed to close database:", err)
		}
	}()

	for i = 0; i < nbUsers; i++ {
		user = User{
			Firstname: fake.Person().FirstName(),
			Lastname:  fake.Person().LastName(),
			Username:  fake.Person().Name(),
			Email:     fake.Internet().Email(),
			Password:  fake.Internet().Password(),
			Roles:     []string{"user"},
		}

		if i == 0 {
			user.Username = testUsername
		}

		DB.Create(&user)
	}

	count, err := CountUsersByUsername(testUsername)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestHashPassword(t *testing.T) {
	var password = "password"

	user := User{
		Firstname: fake.Person().FirstName(),
		Lastname:  fake.Person().LastName(),
		Username:  fake.Person().Name(),
		Email:     fake.Internet().Email(),
		Roles:     []string{"user"},
	}

	if err := user.HashPassword(password); err != nil {
		t.Error(err)
		return
	}

	if user.Password == password {
		t.Error("Expected password to be hashed")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		t.Error("Expected to match password")
		return
	}
}

func TestComparePassword(t *testing.T) {
	var password = "password"

	user := User{
		Firstname: fake.Person().FirstName(),
		Lastname:  fake.Person().LastName(),
		Username:  fake.Person().Name(),
		Email:     fake.Internet().Email(),
		Password:  password,
		Roles:     []string{"user"},
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		t.Error(err)
		return
	}

	user.Password = string(bytes)

	if !user.ComparePassword(password) {
		t.Error("Expected password to match")
	}

	if user.ComparePassword("wrongpassword") {
		t.Error("Expected password to not match")
	}
}

func TestIsRole(t *testing.T) {
	user := User{
		Firstname: fake.Person().FirstName(),
		Lastname:  fake.Person().LastName(),
		Username:  fake.Person().Name(),
		Email:     fake.Internet().Email(),
		Password:  fake.Internet().Password(),
		Roles:     []string{"user"},
	}

	if !user.IsRole("user") {
		t.Error("Expected user to have role 'user'")
	}

	if user.IsRole("admin") {
		t.Error("Expected user to not have role 'admin'")
	}
}
