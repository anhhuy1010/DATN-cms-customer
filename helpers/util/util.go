package util

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/smtp"
	"reflect"
	"strings"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var vnLocation, _ = time.LoadLocation("Asia/Ho_Chi_Minh")

func GenerateUUID() string {
	uuidNew, _ := uuid.NewUUID()
	return uuidNew.String()
}

func ShoudBindHeader(c *gin.Context) bool {
	platform := c.Request.Header.Get("X-PLATFORM")
	lang := c.Request.Header.Get("X-LANG")
	return platform != "" && lang != ""
}

// Lấy thời gian hiện tại theo múi giờ Việt Nam
func NowVN() time.Time {
	return time.Now().In(vnLocation)
}

// Chuyển bất kỳ thời gian nào sang múi giờ Việt Nam
func ToVN(t time.Time) time.Time {
	return t.In(vnLocation)
}

// Debug in ra JSON có format đẹp
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
	UserName string     `json:"username"`
	Email    string     `json:"email"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims
func (c *Claims) Valid() error {
	return nil
}

// GenerateJWT tạo token có thể chứa StartDay và EndDay là nil
func GenerateJWT(uuid, username, email string, startday, endday *time.Time) (string, error) {
	cfg := config.GetConfig()
	jwtKeyStr := cfg.GetString("auth.key")
	jwtKey := []byte(jwtKeyStr)

	now := NowVN()
	claims := &Claims{
		Uuid:     uuid,
		UserName: username,
		Email:    email,
		StartDay: startday,
		EndDay:   endday,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // Token hết hạn sau 24h
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Tạo OTP ngẫu nhiên 6 chữ số
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Gửi mã OTP đến email
func SendOTPEmail(toEmail string, otp string) error {
	from := "tranbaoanhhuy6@gmail.com"
	password := "nrwz hoxd tfvs gldy" // Chú ý: Không nên commit mật khẩu thật vào mã nguồn
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf("Subject: Xác minh tài khoản\n\nMã OTP của bạn là: %s", otp))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		return err
	}

	return nil
}
