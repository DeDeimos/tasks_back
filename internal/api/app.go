package app

import (
	"awesomeProject/internal/app/config"
	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/dsn"
	"awesomeProject/internal/app/redis"
	"awesomeProject/internal/app/repository"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type Application struct {
	repository *repository.Repository
	config     *config.Config
	redis      *redis.Client
	token      *config.TokenConfig
}

func New(ctx context.Context) (*Application, error) {

	_ = godotenv.Load()
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	log.Println("repo done")
	conf, err := config.NewConfig()
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	log.Println("conf done", conf, conf.Redis)

	redisClient, err := redis.New(conf.Redis)
	if err != nil {
		return nil, err
	}
	log.Println("redis done")

	token, err := config.New()
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	return &Application{repository: repo, config: conf, redis: redisClient, token: token}, nil
}

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Role        string `json:"role"`
	Name string `json:"name"`
}

// @Summary Login
// @Description Login
// @Tags auth
// @ID login
// @Accept json
// @Produce json
// @Param input body loginReq true "login info"
// @Success 200 {object} loginResp
// @Router /auth/login [post]
func (a *Application) Login(gCtx *gin.Context) {
	req := &loginReq{}
	log.Println(a)
	cfg := a.token
	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Println(req)

	user, err := a.repository.FindByLogin(req.Login)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Println(user)

	if req.Login == user.Email && generateHashString(req.Password) == user.Password {
		// значит проверка пройдена
		cfg.JWT.SigningMethod = jwt.SigningMethodHS256
		cfg.JWT.ExpiresIn = time.Hour
		// генерируем ему jwt
		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "admin",
			},
			User_ID: user.User_id, // test uuid
			Scopes:  []string{},   // test data
			Role:    user.Role,
			Name: user.Name,
		})
		log.Println(token)
		log.Println(token.Claims.(*ds.JWTClaims).User_ID)
		if token == nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
			return
		}

		gCtx.JSON(http.StatusOK, loginResp{
			ExpiresIn:   int(a.token.JWT.ExpiresIn),
			AccessToken: strToken,
			TokenType:   "Bearer",
			Role:        user.Role,
			Name: user.Name,
		})
	}

	gCtx.AbortWithStatus(http.StatusForbidden) // отдаем 403 ответ в знак того что доступ запрещен
}

type registerReq struct {
	Email string `json:"email"` // лучше назвать то же самое что login
	Name  string `json:"name"`
	Pass  string `json:"pass"`
}

type registerResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Role        string `json:"role"`
}

// @Summary Registration
// @Description Registration
// @Tags auth
// @ID registration
// @Accept json
// @Produce json
// @Param input body registerReq true "user info"
// @Success 200 {object} registerResp
// @Router /auth/registration [post]
func (a *Application) Register(gCtx *gin.Context) {
	req := &registerReq{}

	err := json.NewDecoder(gCtx.Request.Body).Decode(req)
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if req.Pass == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.Name == "" {
		gCtx.AbortWithError(http.StatusBadRequest, fmt.Errorf("name is empty"))
		return
	}

	err = a.repository.CreateUser(ds.User{
		Role:     "user",
		Email:    req.Name,
		Password: generateHashString(req.Pass), // пароли делаем в хешированном виде и далее будем сравнивать хеши, чтобы их не угнали с базой вместе
	})
	if err == nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user, error := a.repository.FindByLogin(req.Email)
	if error != nil {

	}
	cfg := a.token
	// генерируем ему jwt
	token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "admin",
		},
		User_ID: user.User_id, // test uuid
		Scopes:  []string{},   // test data
		Role:    user.Role,
	})
	log.Println(token)
	if token == nil {
		gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
		return
	}

	strToken, err := token.SignedString([]byte(cfg.JWT.Token))
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
		return
	}

	gCtx.JSON(http.StatusOK, loginResp{
		ExpiresIn:   int(a.token.JWT.ExpiresIn),
		AccessToken: strToken,
		TokenType:   "Bearer",
		Role:        user.Role,
	})
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// @Summary Logout
// @Security ApiKeyAuth
// @Description Logout
// @Tags auth
// @ID logout
// @Produce json
// @Success 200 {string} string
// @Router /auth/logout [get]
func (a *Application) Logout(gCtx *gin.Context) {
	// получаем заголовок
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
		gCtx.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа

		return // завершаем обработку
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.token.JWT.Token), nil
	})
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)

		return
	}

	// сохраняем в блеклист редиса
	err = a.redis.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, a.token.JWT.ExpiresIn)
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	gCtx.Status(http.StatusOK)
}
