package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/muhsufyan/api-mongodb/database"
	helper "github.com/muhsufyan/api-mongodb/helpers"
	"github.com/muhsufyan/api-mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// mengakses collection user (mengakses tabel user)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// bandingkan password di db dg password yg diinput apakah sama
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("password is wrong")
		check = false
	}
	return check, msg
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// tangkap data untuk buat data user baru ke db
		var data models.User
		// binding data dlm bntk json
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// validasi data
		validationErr := validate.Struct(data)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// cek jika email telah digunakan
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": data.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured when checking the email"})
		}
		// hash password
		password := HashPassword(*data.Password)
		data.Password = &password
		// cek jika phone telah digunakan
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": data.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured when checking the phone"})
		}
		// cek email/phone/keduanya tlh digunakan
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or phone already exist"})
		}
		// mengisi data create, update dan id
		data.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		data.Update_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		data.ID = primitive.NewObjectID()
		data.User_id = data.ID.Hex()
		// create token dan dlm token payloadnya email, first name, last name, user type dan user id
		token, refreshToken, _ := helper.GenerateAllTokens(*data.Email, *data.First_name, *data.Last_name, *data.User_type, *&data.User_id)
		// set token
		data.Token = &token
		data.Refresh_token = &refreshToken
		// insert data diatas (simpan) kedlm db
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, data)
		if insertErr != nil {
			message := fmt.Sprintf("user data not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": message})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// tangkap data untuk login
		var data models.User
		var foundUser models.User

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// cari email login di db
		err := userCollection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password invalid"})
			return
		}
		// verifikasi password
		passwordIsValid, msg := VerifyPassword(*data.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// tampilkan perpage/ paginasi
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			recordPerPage = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		// cek clientnya admin/user
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// tangkap data dg struct User pd models/UserModel.go
		var data models.User
		// cari data user dg id di db
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&data)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
