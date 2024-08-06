package fixtures

import (
	"fmt"
	"golang-api/models"
	"math/rand"
)

func LoadUsers() error {
	if err := models.ClearUsers(); err != nil {
		return err
	}

	for i := 0; i < USER_NB; i++ {
		role := USER_ROLES[rand.Intn(len(USER_ROLES))]

		firstName := fake.Person().FirstName()
		lastName := fake.Person().LastName()
		defaultAvatarUrl := fmt.Sprintf("https://api.dicebear.com/9.x/initials/svg?seed=%s%s", firstName, lastName)

		user := models.User{
			Firstname:       firstName,
			Lastname:        lastName,
			Username:        fake.Person().FirstName(),
			Email:           fake.Internet().Email(),
			ProfileImageUrl: defaultAvatarUrl,
			Roles:           []string{role},
			Verified:        true,
		}

		if err := user.HashPassword(USER_PASSWORD); err != nil {
			return err
		}

		if err := user.Save(); err != nil {
			return err
		}

		fmt.Printf("User %s created with password %s\n", user.Username, USER_PASSWORD)
	}

	return nil
}
