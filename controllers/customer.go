package controllers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/anhhuy1010/DATN-cms-customer/constant"
	"github.com/anhhuy1010/DATN-cms-customer/grpc"
	pbUsers "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/users"
	pbIdeas "github.com/anhhuy1010/DATN-cms-ideas/grpc/proto/idea"

	"github.com/anhhuy1010/DATN-cms-customer/helpers/respond"
	"github.com/anhhuy1010/DATN-cms-customer/helpers/util"
	"github.com/anhhuy1010/DATN-cms-customer/middleware"
	"github.com/anhhuy1010/DATN-cms-customer/models"
	request "github.com/anhhuy1010/DATN-cms-customer/request/user"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
}

func (userCtl UserController) SignUp(c *gin.Context) {
	var req request.SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	customerSignup := models.Customer{}
	customerSignup.Uuid = util.GenerateUUID()
	customerSignup.UserName = req.UserName
	customerSignup.Password = req.Password
	customerSignup.Email = req.Email
	customerSignup.IsActive = 0 // Mặc định inactive, chờ xác thực OTP
	customerSignup.StartDay = nil
	customerSignup.EndDay = nil

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customerSignup.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Invalid password"))
		return
	}
	customerSignup.Password = string(hashedPassword)

	_, err = customerSignup.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, respond.UpdatedFail())
		return
	}

	// ---- Generate OTP ----
	otpCode := util.GenerateOTP()

	// Gửi OTP tới email
	err = util.SendOTPEmail(customerSignup.Email, otpCode)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể gửi mã OTP tới email"})
		return
	}

	// ---- Lưu OTP vào database ----
	otp := models.OTP{
		Uuid:      util.GenerateUUID(),
		UserUuid:  customerSignup.Uuid,
		OtpCode:   otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute), // OTP có hiệu lực 5 phút
		CreatedAt: time.Now(),
	}

	err = otp.Insert(c.Request.Context())
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng ký thành công, vui lòng kiểm tra email để xác thực tài khoản!",
		"uuid":    customerSignup.Uuid,
	})
}

func (userCtl UserController) VerifyOTP(c *gin.Context) {
	OTPModel := models.OTP{}
	userModel := models.Customer{}
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Tìm OTP theo Email
	otpRecord, err := OTPModel.FindOTPByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP không tồn tại hoặc đã hết hạn"})
		return
	}

	// So sánh OTP
	if otpRecord.OtpCode != req.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã OTP không đúng"})
		return
	}

	// Update tài khoản: is_active = 1
	customer, err := userModel.FindCustomerByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy tài khoản"})
		return
	}
	customer.IsActive = 1
	_, err = customer.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Xác thực thất bại"})
		return
	}
	c.JSON(http.StatusOK, respond.Success(customer.Uuid, "update successfully"))

	// Xoá OTP
	_ = OTPModel.DeleteOTP(context.Background(), req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Xác thực thành công!"})
}

func (userCtl UserController) Login(c *gin.Context) {
	userModel := models.Customer{}

	var req request.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"email": req.Email}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("user not found"))
		return
	}

	token, err := util.GenerateJWT(user.Uuid, user.StartDay, user.EndDay)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("create token found"))
		return
	}
	userLogin := models.Tokens{}
	userLogin.UserUuid = user.Uuid
	userLogin.Uuid = util.GenerateUUID()
	userLogin.Token = token
	userLogin.IsDelete = 0

	_, err = userLogin.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(request.LoginResponse{Token: token}, "login successfully"))
}

//////////////////////////////////////////////////////////////////////

func (userCtl UserController) Logout(c *gin.Context) {
	// Lấy token từ header
	tokenStr := c.GetHeader("x-token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Xóa token khỏi CSDL
	tokens := models.Tokens{}
	condition := bson.M{"token": tokenStr, "is_delete": 0}
	update := bson.M{"$set": bson.M{"is_delete": 1}}

	err := tokens.UpdateOne(condition, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

////////////////////////////////////////////////////////////////////////////

func (userCtl UserController) CheckRole(token string) (*pbUsers.DetailResponse, error) {
	grpcConn := grpc.GetInstance()
	client := pbUsers.NewUserClient(grpcConn.UsersConnect)
	req := pbUsers.DetailRequest{
		Token: token,
	}
	resp, err := client.Detail(context.Background(), &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (userCtl UserController) GetRoleByToken(token string) (*request.CheckRoleResponse, error) {
	tokenModel := models.Tokens{}
	userModel := models.Customer{}

	condition := bson.M{"token": token}
	tokenDoc, err := tokenModel.FindOne(condition)
	if err != nil {
		return nil, errors.New("token not found")
	}
	if tokenDoc == nil {
		return nil, errors.New("token document is nil")
	}

	cond := bson.M{"uuid": tokenDoc.UserUuid}
	user, err := userModel.FindOne(cond)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := &request.CheckRoleResponse{
		UserUuid: user.Uuid,
	}
	return resp, nil
}

func RoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		userCtl := UserController{}
		resp, err := userCtl.GetRoleByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("startday", resp.StartDay)
		c.Set("endday", resp.EndDay)
		c.Set("user_uuid", resp.UserUuid)
		c.Next()
	}
}
func (userCtl UserController) List(c *gin.Context) {
	userModel := new(models.Customer)
	var req request.GetListRequest
	err := c.ShouldBindWith(&req, binding.Query)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	cond := bson.M{}
	if req.Username != nil {
		cond["username"] = req.Username
	}

	if req.IsActive != nil {
		cond["is_active"] = req.IsActive
	}
	if req.Role != nil {
		cond["role"] = req.Role
	}

	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)
	var respData []request.ListResponse
	users, err := userModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	for _, user := range users {
		res := request.ListResponse{
			Uuid:     user.Uuid,
			IsActive: user.IsActive,
		}
		respData = append(respData, res)
	}
	total, err := userModel.Count(c, cond)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Detail(c *gin.Context) {
	userModel := new(models.Customer)
	var reqUri request.GetDetailUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	response := request.GetDetailResponse{
		Uuid:     user.Uuid,
		Email:    user.Email,
		IsActive: user.IsActive,
		IsDelete: user.IsDelete,
	}

	c.JSON(http.StatusOK, respond.Success(response, "Successfully"))
}

// ////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Update(c *gin.Context) {
	userModel := new(models.Customer)
	var reqUri request.UpdateUri

	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.UpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Delete(c *gin.Context) {
	userModel := new(models.Customer)
	var reqUri request.DeleteUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsDelete = constant.DELETE

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "Delete successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) UpdateStatus(c *gin.Context) {
	userModel := new(models.Customer)
	var reqUri request.UpdateStatusUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	var req request.UpdateStatusRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if *req.IsActive < 0 || *req.IsActive > 1 {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Stauts just can be set in range [0..1]"))
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User no found!"))
		return
	}

	user.IsActive = *req.IsActive

	_, err = user.Update()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) Create(c *gin.Context) {
	var req request.GetInsertRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}
	userData := models.Customer{}
	userData.Uuid = util.GenerateUUID()

	userData.Password = req.Password
	userData.Email = req.Email
	userData.IsActive = 1
	userData.Password = req.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
		return
	}
	userData.Password = string(hashedPassword)

	_, err = userData.Insert()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}
	c.JSON(http.StatusOK, respond.Success(userData.Uuid, "create successfully"))
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) PostIdea(ctx *gin.Context) {
	var req struct {
		IdeasName     string `json:"ideas_name"`
		Industry      string `json:"industry"`
		ContentDetail string `json:"content_detail"`
		Price         int32  `json:"price"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Lấy claims từ middleware
	claims, ok := middleware.ExtractClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Không xác thực được người dùng"})
		return
	}
	customerUuid, _ := claims["uuid"].(string)
	customerName, _ := claims["name"].(string)
	customerEmail, _ := claims["email"].(string)

	if customerUuid == "" || customerName == "" || customerEmail == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu thông tin người dùng trong token"})
		return
	}

	// Gọi gRPC
	grpcConn := grpc.GetInstance()
	client := pbIdeas.NewIdeaServiceClient(grpcConn.IdeasConnect)

	res, err := client.CreateIdea(context.Background(), &pbIdeas.CreateIdeaRequest{
		IdeasName:     req.IdeasName,
		Industry:      req.Industry,
		ContentDetail: req.ContentDetail,
		Price:         req.Price,
		CustomerUuid:  customerUuid,
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": res.Message,
		"uuid":    res.Uuid,
	})
}

// //////////////////////////////////////////////////////////////////

func AddFavorite(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	customerUuid := c.GetString("customer_uuid")
	if customerUuid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chưa đăng nhập"})
		return
	}

	fav := models.Favorite{
		Uuid:         util.GenerateUUID(),
		CustomerUuid: customerUuid,
		PostUuid:     req.PostUuid,
		PostType:     req.PostType,
		CreatedAt:    time.Now(),
	}

	_, err := fav.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể thêm vào danh sách yêu thích"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã thêm vào danh sách yêu thích"})
}
