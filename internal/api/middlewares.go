package app

import (
	"awesomeProject/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

const jwtPrefix = "Bearer "

func (a *Application) WithAuthCheck(assignedRoles ...string) func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		log.Println(jwtStr)
		log.Println(1)
		log.Println(assignedRoles)
		if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
			gCtx.AbortWithStatus(http.StatusForbidden) // отдаем что нет доступа

			return // завершаем обработку
		}
		log.Println(11)
		log.Println(assignedRoles)
		// отрезаем префикс
		jwtStr = jwtStr[len(jwtPrefix):]

		err := a.redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)

			return
		}
		log.Println(12)
		log.Println(assignedRoles)
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			gCtx.AbortWithError(http.StatusInternalServerError, err)

			return
		}
		log.Println(jwtStr)
		log.Println(2)
		log.Println(assignedRoles)

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.token.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}
		log.Println(3)
		log.Println(assignedRoles)
		log.Println(token)
		myClaims := token.Claims.(*ds.JWTClaims)
		ctxWithUserID := gCtx.Request.Context()
		ctxWithUserID = context.WithValue(ctxWithUserID, "userID", myClaims.User_ID)
		gCtx.Set("userID", myClaims.User_ID)

		userID, exists := gCtx.Get("userID")
		if exists {
			fmt.Println(userID.(uint))
		}
		log.Println(4)
		log.Println(assignedRoles)
		ctxWithUserRole := gCtx.Request.Context()
		ctxWithUserRole = context.WithValue(ctxWithUserRole, "userRole", myClaims.Role)
		gCtx.Set("userRole", myClaims.Role)
		log.Println(5)
		log.Println(assignedRoles)
		userRole, exists := gCtx.Get("userRole")
		if exists {
			fmt.Println(userRole.(string))
		}
		log.Println(6)
		log.Println(assignedRoles)
		fmt.Println("Сюда()")
		fmt.Println(myClaims)
		authorized := false

		for _, oneOfAssignedRole := range assignedRoles {
			if myClaims.Role == oneOfAssignedRole {
				authorized = true
				break
			}
		}

		if !authorized {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Printf("role %s is not assigned in %s", myClaims.Role, assignedRoles)
			return
		}

	}

}

// мне нужен midleware который будет проверять наличие токена и если он есть в контекст класть userID и userRole
// если токена нет то в контекст кладем пустые строки
func (a *Application) WithOptionalCheck() func(ctx *gin.Context) {
	return func(gCtx *gin.Context) {
		jwtStr := gCtx.GetHeader("Authorization")
		log.Println(jwtStr)
		log.Println(1)

		if jwtStr == "" {
			// отдаем что нет доступа
			gCtx.Set("userID", "")
			gCtx.Set("userRole", "")
			gCtx.Next()
			return
		}
		log.Println(11)
		// отрезаем префикс
		jwtStr = jwtStr[len(jwtPrefix):]

		err := a.redis.CheckJWTInBlacklist(gCtx.Request.Context(), jwtStr)
		if err == nil { // значит что токен в блеклисте
			gCtx.AbortWithStatus(http.StatusForbidden)

			return
		}
		log.Println(12)
		if !errors.Is(err, redis.Nil) { // значит что это не ошибка отсуствия - внутренняя ошибка
			gCtx.AbortWithError(http.StatusInternalServerError, err)

			return
		}
		log.Println(jwtStr)
		log.Println(2)

		token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.token.JWT.Token), nil
		})
		if err != nil {
			gCtx.AbortWithStatus(http.StatusForbidden)
			log.Println(err)

			return
		}
		log.Println(3)
		log.Println(token)
		myClaims := token.Claims.(*ds.JWTClaims)
		ctxWithUserID := gCtx.Request.Context()
		ctxWithUserID = context.WithValue(ctxWithUserID, "userID", myClaims.User_ID)
		gCtx.Set("userID", myClaims.User_ID)

		userID, exists := gCtx.Get("userID")
		if exists {
			fmt.Println(userID.(uint))
		}
		log.Println(4)
		ctxWithUserRole := gCtx.Request.Context()
		ctxWithUserRole = context.WithValue(ctxWithUserRole, "userRole", myClaims.Role)
		gCtx.Set("userRole", myClaims.Role)
		log.Println(5)

		userRole, exists := gCtx.Get("userRole")
		if exists {
			fmt.Println(userRole.(string))
		}
		log.Println(6)
		fmt.Println("Сюда()")
		fmt.Println(myClaims)

		fmt.Println(myClaims.User_ID)

		gCtx.Next()
	}

}
