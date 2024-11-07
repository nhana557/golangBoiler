package controller

import (
	"boiler-go/entities"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type LoginHandler struct {
	LoginUsecase entities.LoginUsecase
	Config       *viper.Viper
}

func NewLoginHandler(e *gin.RouterGroup, lu entities.LoginUsecase, config *viper.Viper) {
	handler := &LoginHandler{
		LoginUsecase: lu,
		Config:       config,
	}
	e.POST("/login/admin", handler.CreateJwtAdmin)
	e.POST("/login", handler.CreateJwtUser)
}

func isRequestValid(m *entities.Login) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (login *LoginHandler) CreateJwtUser(c *gin.Context) {
	var (
		err          error
		token        string
		loginPayload entities.Login
	)
	if err = c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ok, err := isRequestValid(&loginPayload); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := login.LoginUsecase.GetUser(ctx, loginPayload.Username, loginPayload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Your username or password were wrong"})
		return
	}

	lifetime, err := strconv.ParseInt(login.Config.GetString("jwt.lifetime"), 10, 64)
	if err != nil {
		lifetime = 60
	}

	secret := login.Config.GetString("jwt.secret")
	token, err = createJwtToken(res.ID.Hex(), "user", lifetime, secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": lifetime,
	})
}

func (login *LoginHandler) CreateJwtAdmin(c *gin.Context) {
	var (
		err          error
		token        string
		loginPayload entities.Login
	)

	if err = c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if ok, err := isRequestValid(&loginPayload); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ctx := c.Request.Context()
	// if ctx == nil {
	// 	ctx = context.Background()
	// }

	adminUsername := login.Config.GetString("admin.username")
	adminPassword := login.Config.GetString("admin.password")

	if loginPayload.Username == adminUsername && loginPayload.Password == adminPassword {
		// create jwt token
		lifetime, err := strconv.ParseInt(login.Config.GetString("jwt.lifetime"), 10, 64)
		if err != nil {
			lifetime = 60
		}

		secret := login.Config.GetString("jwt.secret")
		token, err = createJwtToken(adminUsername, "admin", lifetime, secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"expires_in": lifetime,
			"token":      token,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Your username or password were wrong"})
}

func createJwtToken(uname, jtype string, lifetime int64, secret string) (string, error) {
	type JwtClaims struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
		jwt.StandardClaims
	}

	var (
		claim    JwtClaims
		lifeTime int64 = time.Now().Add(time.Duration(lifetime) * time.Minute).Unix()
	)

	if jtype == "admin" {
		claim = JwtClaims{
			Name:    uname,
			IsAdmin: true,
			StandardClaims: jwt.StandardClaims{
				Id:        uname,
				ExpiresAt: lifeTime,
			},
		}
	} else {
		claim = JwtClaims{
			Name:    uname,
			IsAdmin: false,
			StandardClaims: jwt.StandardClaims{
				Id:        uname,
				ExpiresAt: lifeTime,
			},
		}
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
