package routes

import (
	"main/config"
	"main/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("secret_key")

func HashPassword(Password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterHandler(ctx *gin.Context){
	items := models.Login{}
	config.DB.Find(&items)
	
	var Register models.Login
	err := ctx.BindJSON(&Register)
	if err != nil {
		errorMessages := []string{"Please fill all requirements"}

		ctx.JSON(404, gin.H{
			"errors" : errorMessages,
		})
		return
	}
	if Register.Username == items.Username {
		ctx.JSON(404, gin.H{"error" : "Username already used"})
		return
	}
	
	/*hash password*/

	// passwordHash := Register.Password
	// hash, _ := HashPassword(passwordHash)
	
	// Register.Password = hash
	// fmt.Println("password", passwordHash)
	// fmt.Println("Hash", hash)

	// match := CheckPasswordHash(passwordHash, hash)
	// fmt.Println(match)
	config.DB.Create(&Register)
	
	ctx.JSON(200, gin.H{
		"Status" : "Succesfully register" ,
	})
}

func LoginHandler(ctx *gin.Context){
	items := models.Login{}
	config.DB.First(&items)
	
	var userData models.Login
	err := ctx.ShouldBindJSON(&userData)
	if err != nil {
		errorMessages := []string{"Please fill all requirements"}
		ctx.JSON(404, gin.H{"errors" : errorMessages,})
		return
	}
	
	if userData.Username != items.Username || userData.Password != items.Password{
		ctx.JSON(404, gin.H{"erro" : "invalid username or password"})
		return
	}
	
	expiredTime := time.Now().Add(time.Minute * 30)

	claims := &Claims{
		Username: userData.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			ctx.JSON(404, gin.H{
				"Status" : "Error to generate token",
				"error"  : err,
			})
		}

	ctx.JSON(200, gin.H{
		"Status" : "Succesfully",
		"Token"  : tokenString, 
	})	
}

func Base(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"Message" : "Berhasil",
	})
}

func GetArticles(ctx *gin.Context){

	items := []models.Articles{}
	config.DB.Find(&items)

	ctx.JSON(200, gin.H{
		"Status" : "Berhasil",
		"Data" : items,
	})
}

func GetArticle(ctx *gin.Context){
	Id := ctx.Param("id")
	var item models.Articles
	
	err := config.DB.First(&item, "id = ?", Id).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"Status" : "Data not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"Status" : "Successfully",
		"Data" : item,
	})
}

func PostArticle(ctx *gin.Context){
	var PostArticle models.Articles
	err := ctx.BindJSON(&PostArticle)
	if err != nil {
		errorMessages := []string{"Please fill all requirements"}

		ctx.JSON(404, gin.H{
			"errors" : errorMessages,
		})
		return
	}
	config.DB.Create(&PostArticle)

	ctx.JSON(200, gin.H{
		"Status" : "Succesfully" ,
		"Data" : PostArticle,
	})
}

func UpdateArticle(ctx *gin.Context){
	Id := ctx.Param("id")
    user := models.Articles{}

	config.DB.First(&user, Id)
		if user.ID == 0 {
			ctx.JSON(404, gin.H{ "Status" : "Not found" })
			return
		}
			if err := ctx.BindJSON(&user); err != nil {
				ctx.JSON(404, err)
				return
			}

    config.DB.Model(&user).Updates(user)
	ctx.JSON(200, gin.H{
		"Status" : "Succesfully",
		"Data" : user,
	})
}

func DeleteArticle(ctx *gin.Context){
	Id := ctx.Param("id")
	var DeleteArticle models.Articles

	config.DB.First(&DeleteArticle, Id)
		if DeleteArticle.ID == 0 {
			ctx.JSON(404, gin.H{ "Status" : "Not found" })
			return
		}

	err := config.DB.Exec("DELETE FROM articles where id = ?", Id).Error
	if err != nil {
		ctx.JSON(404, gin.H {"Error": err})
		return
	}

	ctx.JSON(200, gin.H{
		"Status" : "Succesfully delete",
		"Data" : DeleteArticle,
	})	 
}