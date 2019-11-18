package migration

import (
	"fmt"
	"messaging/messaging"
	"messaging/models"
	"strings"
	"time"
)

//TODO: migration for notification

func GetMessageSeedSuccess() (mod models.Messaging) {
	layout := "2019-10-23 04:40:08.034983+00"
	str := "2019-10-23 04:40:08.034983+00"
	t, _ := time.Parse(layout, str)
	return models.Messaging{
		ClientID:     2,
		Partner:      "adsmedia",
		PhoneNumber:  "081134567892",
		Message:      "kode otp anda 126236",
		Status:       "success",
		ErrorMessage: "",
		SendTime:     t,
		RawResponse:  "{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":10,\"sendingstatustext\":\"success\",\"price\":320}}]}]}",
	}
}
func GetMessageSeedFailed() (mod models.Messaging) {
	layout := "2019-10-23 04:40:08.034983+00"
	str := "2019-10-23 04:40:08.034983+00"
	t, _ := time.Parse(layout, str)
	return models.Messaging{
		ClientID:     2,
		Partner:      "adsmedia",
		PhoneNumber:  "081134567892",
		Message:      "kode otp anda 126236",
		Status:       "failed",
		ErrorMessage: "Invalid Number",
		SendTime:     t,
		RawResponse:  "{\"sending_respon\":[{\"globalstatus\":10,\"globalstatustext\":\"Success\",\"datapacket\":[{\"packet\":{\"number\":\"6282297335657\",\"sendingid\":1287265,\"sendingstatus\":60,\"sendingstatustext\":\"Invalid Number\",\"price\":320}}]}]}",
	}
}

func GetUsersSeed() (mod []models.Users) {
	return []models.Users{
		models.Users{
			Username: "adminkey",
			Password: "adminsecret",
		},
		// models.Users{
		// 	Username: "smsotp",
		// 	Password: "P@ssw0rd",
		// 	Email:    "smsotp@ayannah.com",
		// },
		// models.Users{
		// 	Username: "lendernotif",
		// 	Password: "P@ssw0rd",
		// 	Email:    "lendernotif@ayannah.com",
		// },
		// models.Users{
		// 	Username: "borrowernotif",
		// 	Password: "P@ssw0rd",
		// 	Email:    "borrowernotif@ayannah.com",
		// },
	}
}
func Seed() {
	if messaging.App.ENV == "development" {
		// seed clients
		Clients := []models.Clients{
			models.Clients{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Clients{
				Name:   "Client",
				Key:    "clientkey",
				Role:   "client",
				Secret: "clientsecret",
			},
		}
		for _, internal := range Clients {
			internal.Create()
		}

		messagings := []models.Messaging{
			GetMessageSeedSuccess(),
			GetMessageSeedFailed(),
			// models.Messaging{
			// 	Partner: "OtherPartner",
			// },
		}
		for _, messaging := range messagings {
			messaging.Create()
		}

		//seed users
		for _, users := range GetUsersSeed() {
			users.Create()
		}

		notifications := []models.Notification{
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-1",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 1 diterima oleh Bank A 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-2",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 2 diterima oleh Bank B 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "lender-1",
				Title:         "Promosi",
				MessageBody:   "Promosi Produk A...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      1,
				RecipientID:   "admin-1",
				Title:         "Info",
				MessageBody:   "Info, informasi...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
		}
		for _, notification := range notifications {
			notification.Create()
		}

	}
}

func TestSeed() {
	if messaging.App.ENV == "development" {
		// seed Clients
		Clients := []models.Clients{
			models.Clients{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Clients{
				Name:   "Client",
				Key:    "clientkey",
				Role:   "client",
				Secret: "clientsecret",
			},
		}
		for _, internal := range Clients {
			internal.Create()
		}

		messagings := []models.Messaging{
			// models.Messaging{
			// 	Partner: "adsmedia",
			// },
			// models.Messaging{
			// 	Partner: "Partner",
			// },
			GetMessageSeedSuccess(),
			GetMessageSeedFailed(),
		}
		for _, messaging := range messagings {
			messaging.Create()
		}

		//seed users :
		// for _, users := range GetUsersSeed() {
		// 	users.Create()
		// }

		notifications := []models.Notification{
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-1",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 1 diterima oleh Bank A 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-2",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 2 diterima oleh Bank B 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "lender-1",
				Title:         "Promosi",
				MessageBody:   "Promosi Produk A...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      1,
				RecipientID:   "admin-1",
				Title:         "Info",
				MessageBody:   "Info, informasi...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
		}
		for _, notification := range notifications {
			notification.Create()
		}

	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"clients",
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
