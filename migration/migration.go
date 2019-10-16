package migration

import (
	"fmt"
	"messaging/messaging"
	"messaging/models"
	"strings"
)

func Seed() {
	if messaging.App.ENV == "development" {
		// seed internals
		internals := []models.Internals{
			models.Internals{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Internals{
				Name:   "Client",
				Key:    "clientkey",
				Role:   "client",
				Secret: "clientsecret",
			},
		}
		for _, internal := range internals {
			internal.Create()
		}

		messagings := []models.Messaging{
			models.Messaging{
				Partner: "adsmedia",
			},
			models.Messaging{
				Partner: "Partner",
			},
		}
		for _, messaging := range messagings {
			messaging.Create()
		}
	}
}

func TestSeed() {
	if messaging.App.ENV == "development" {
		// seed internals
		internals := []models.Internals{
			models.Internals{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Internals{
				Name:   "Client",
				Key:    "clientkey",
				Role:   "client",
				Secret: "clientsecret",
			},
		}
		for _, internal := range internals {
			internal.Create()
		}

		messagings := []models.Messaging{
			models.Messaging{
				Partner: "adsmedia",
			},
			models.Messaging{
				Partner: "Partner",
			},
		}
		for _, messaging := range messagings {
			messaging.Create()
		}

	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"internals",
				"messagings",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = messaging.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
