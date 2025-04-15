package util

import (
	"context"
	"encoding/json"

	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateUUID() (s string) {
	uuidNew, _ := uuid.NewUUID()
	return uuidNew.String()
}

func ShoudBindHeader(c *gin.Context) bool {
	platform := c.Request.Header.Get("X-PLATFORM")
	lang := c.Request.Header.Get("X-LANG")

	if platform == "" || lang == "" {
		return false
	}

	return true
}

func GetNowUTC() time.Time {
	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	return currentTime
}
func DebugJson(value interface{}) {
	fmt.Println(reflect.TypeOf(value).String())
	prettyJSON, _ := json.MarshalIndent(value, "", "    ")
	fmt.Printf("%s\n", string(prettyJSON))
}

func GetKeyFromContext(ctx context.Context, key string) (interface{}, bool) {
	if v := ctx.Value(key); v != nil {
		return v, true
	}

	return nil, false
}

func LogPrint(jsonData interface{}) {
	prettyJSON, _ := json.MarshalIndent(jsonData, "", "")
	fmt.Printf("%s\n", strings.ReplaceAll(string(prettyJSON), "\n", ""))
}

type Claims struct {
	Uuid     string     `json:"uuid"`
	StartDay *time.Time `json:"startday"`
	EndDay   *time.Time `json:"endday"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims (trả về nil để chấp nhận mọi token hợp lệ).
func (c *Claims) Valid() error {
	return nil
}

// GenerateJWT tạo token có thể chứa StartDay và EndDay là nil
func GenerateJWT(uuid string, startday, endday *time.Time) (string, error) {
	cfg := config.GetConfig()
	jwtKeyStr := cfg.GetString("auth.key")
	jwtKey := []byte(jwtKeyStr)

	claims := &Claims{
		Uuid:     uuid,
		StartDay: startday,
		EndDay:   endday,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // hoặc nil nếu không muốn expire
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
