package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nkyizbay/shop-project/internal/auth"
)

type handler struct {
	userService  Service
	JwtSecretKey string
}

func NewHandler(rout *gin.Engine, userService Service, jwtSecretKey string) *handler {
	h := handler{
		userService:  userService,
		JwtSecretKey: jwtSecretKey,
	}

	rout.POST("/register", h.Register)
	rout.POST("/login", h.Login)
	rout.GET("/logout", h.Logout)

	return &h
}

func (h *handler) Register(c *gin.Context) {
	user := new(User)

	if err := c.Bind(&user); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if user.IsNameEmpty() {
		c.String(http.StatusBadRequest, "Empty username warning")
		return
	}

	if user.IsUserTypeInvalid() {
		c.String(http.StatusBadRequest, "Invalid user type warning")
		return
	}

	if user.IsAuthTypeInvalid() {
		c.String(http.StatusBadRequest, "Invalid auth type warning")
		return
	}

	if user.IsPasswordInvalid() {
		c.String(http.StatusBadRequest, "Password should be between 5 and 12 characters")
		return
	}

	password, err := user.HashPassword()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	user.Password = password

	requestCtx := c.Request.Context()

	err = h.userService.Register(requestCtx, user)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
		return
	}

	c.String(http.StatusCreated, "User created")
}

func (h *handler) Login(c *gin.Context) {
	var credentials auth.Credentials
	if err := c.Bind(&credentials); err != nil {
		c.String(http.StatusBadRequest, "non valid credentials")
		return
	}

	user, err := h.userService.Login(c.Request.Context(), credentials)
	if err != nil {
		c.String(http.StatusInternalServerError, "error internal")
		return
	}

	expirationTime := &jwt.NumericDate{Time: time.Now().Add(time.Hour)}
	claims := auth.Claims{
		Username: user.UserName,
		UserType: user.UserType,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.JwtSecretKey))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error, something went wrong")
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime.Time
	http.SetCookie(c.Writer, cookie)

	c.String(http.StatusOK, "successful log")
}

func (h *handler) Logout(c *gin.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.MaxAge = 0
	http.SetCookie(c.Writer, cookie)

	c.String(http.StatusOK, "You have successfully logout")
}
