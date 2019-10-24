package handlers

import (
	"fmt"
	"messaging/messaging"
	"messaging/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type (
	AdminLoginCreds struct {
		Username string `json:"key"`
		Password string `json:"password"`
	}
)

// admin login
func AdminLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		credentials AdminLoginCreds
		user        models.Users
		validKey    bool
		token       string
		err         error
	)

	rules := govalidator.MapData{
		"username": []string{"required"},
		"password": []string{"required"},
	}

	validate := validateRequestPayload(c, rules, &credentials)
	if validate != nil {
		return returnInvalidResponse(http.StatusBadRequest, validate, "invalid login")
	}

	// check if theres record
	validKey = messaging.App.DB.Where("username = ?", credentials.Username).Find(&user).RecordNotFound()

	if !validKey { // check the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			return returnInvalidResponse(http.StatusUnauthorized, err, "invalid login")
		}

		// if user.Status == false {
		// 	return returnInvalidResponse(http.StatusUnauthorized, err, "invalid login")
		// }

		//TODO: cek this out
		tokenrole := "admin"
		token, err = createJwtToken(strconv.FormatUint(user.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
	} else {
		return returnInvalidResponse(http.StatusUnauthorized, "", "invalid login")
	}

	jwtConf := messaging.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", messaging.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute
	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
