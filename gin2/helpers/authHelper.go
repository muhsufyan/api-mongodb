package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// cek tipe clientnya
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("your access is denied for access this resource")
	}
	return err
}

// cek user id
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	// client adlh user dan hanya dpt mengakses data dirinya sendiri(tdk dpt mengakses data org lain)
	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
