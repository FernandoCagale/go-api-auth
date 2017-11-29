package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/FernandoCagale/go-api-auth/src/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Save(c echo.Context) error {
	user := new(models.User)

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "BadRequest",
		})
	}

	if errors, valid := user.Validate(); !valid {
		return c.JSON(http.StatusBadRequest, errors)
	}

	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users := []models.User{}

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := db.Find(&users).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Get(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := db.Find(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "NotFound",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	userBind := new(models.User)
	user := new(models.User)

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := c.Bind(userBind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "BadRequest",
		})
	}

	if errors, valid := userBind.Validate(); !valid {
		return c.JSON(http.StatusBadRequest, errors)
	}

	if err := db.Find(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "NotFound",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := db.Model(&user).Updates(&userBind).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := db.Find(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "NotFound",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Deleted",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	user := new(models.User)
	userBind := new(models.User)

	db, valid := getConnection(c)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "InternalServerError",
		})
	}

	if err := c.Bind(userBind); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "BadRequest",
		})
	}

	if err := db.Where("username = ?", userBind.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "StatusUnauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	if valid := user.ValidatePassword(userBind.Password); !valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "StatusUnauthorized",
		})
	}

	token, err := createJwtToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "InternalServerError",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func createJwtToken(user *models.User) (string, error) {
	claims := models.JwtClaims{
		user.Username,
		jwt.StandardClaims{
			Id:        strconv.Itoa(user.Id),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return rawToken.SignedString([]byte("secret"))
}

func getConnection(c echo.Context) (*gorm.DB, bool) {
	db := c.Get("db")
	if db != nil {
		return db.(*gorm.DB), true
	}
	return nil, false
}
