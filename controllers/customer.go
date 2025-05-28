package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	userModel := models.Customer{}
	var req request.SignUpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"email": req.Email}
	_, err = userModel.FindOne(condition)
	if err == nil {
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("Tài khoản đã được đăng ký"))
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
	customerSignup.Image = ""
	customerSignup.Introduce = ""
	customerSignup.IsDelete = 0

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
		Email:     customerSignup.Email,
		OtpCode:   otpCode,
		ExpiresAt: util.NowVN().Add(5 * time.Minute), // OTP có hiệu lực 5 phút
		CreatedAt: util.NowVN(),
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
		OtpCode string `json:"otp_code"`
		Email   string `json:"email"`
	}

	// Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng nhập mã OTP"})
		return
	}

	// Tìm OTP theo Email
	otpRecord, err := OTPModel.FindOTPByEmail(context.Background(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP không tồn tại hoặc đã hết hạn"})
		return
	}
	if time.Now().After(otpRecord.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã OTP đã hết hạn"})
		_ = OTPModel.DeleteOTP(context.Background(), otpRecord.OtpCode) // xóa nếu hết hạn
		return
	}

	// Log để debug
	log.Printf("DEBUG: From DB - Email=%s, OTP=%s", otpRecord.Email, otpRecord.OtpCode)
	log.Printf("DEBUG: From request - Email=%s, OTP=%s", req.Email, req.OtpCode)

	// So sánh OTP
	if otpRecord.OtpCode != req.OtpCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã OTP không đúng"})
		return
	}

	// Cập nhật tài khoản: is_active = 1
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

	// Xoá OTP (nếu thất bại vẫn tiếp tục)
	if err := OTPModel.DeleteOTP(context.Background(), req.OtpCode); err != nil {
		log.Println("Lỗi khi xoá OTP:", err)
	}

	// Trả về thành công
	c.JSON(http.StatusOK, gin.H{
		"uuid":    customer.Uuid,
		"message": "Xác thực thành công!",
	})
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
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("không tìm thấy người dùng!"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("wrong password"))
		return
	}

	token, err := util.GenerateJWT(user.Uuid, user.UserName, user.Email, user.StartDay, user.EndDay)
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
	userLogin.UserEmail = user.Email
	userLogin.UserName = user.UserName

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
		UserName: user.UserName,
		Email:    user.Email,
		StartDay: user.StartDay,
		EndDay:   user.EndDay,
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

		// ✅ Tạo claims từ resp
		claims := &util.Claims{
			Uuid:     resp.UserUuid,
			UserName: resp.UserName,
			Email:    resp.Email,
			StartDay: resp.StartDay,
			EndDay:   resp.EndDay,
		}

		// ✅ Gán vào context để ExtractClaims dùng được
		c.Set("claims", claims)
		c.Set("customer_uuid", claims.Uuid)
		c.Set("customer_name", claims.UserName)

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
			UserName: user.UserName,
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

	customerUuid, exists := c.Get("customer_uuid")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid is missing"})
		return
	}
	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}

	var reqUri request.UpdateUri
	// Validation input
	err := c.ShouldBindUri(&reqUri)
	if err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	if customerUuidStr != reqUri.Uuid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this user"})
		return
	}

	var req request.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	condition := bson.M{"uuid": reqUri.Uuid}
	user, err := userModel.FindOne(condition)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.ErrorCommon("User not found!"))
		return
	}

	if req.UserName != "" {
		user.UserName = req.UserName
	}
	if req.Image != "" {
		user.Image = req.Image
		fmt.Println("---------------------------Image from request:", req.Image)
	}
	if req.Introduce != "" {
		user.Introduce = req.Introduce
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, respond.ErrorCommon("invalid password"))
			return
		}
		user.Password = string(hashedPassword)
	}

	if _, err := user.Update(); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, respond.UpdatedFail())
		return
	}

	c.JSON(http.StatusOK, respond.Success(user.Uuid, "update successfully"))
}

// ///////////////////////////////////////////////////////////////////////////
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
		IdeasName      string `json:"ideas_name"`
		Industry       string `json:"industry"`
		ContentDetail  string `json:"content_detail"`
		Price          int32  `json:"price"`
		Image          string `json:"image"` // Client gửi base64 trực tiếp
		Procedure      string `json:"is_procedure"`
		Value_Benefits string `json:"value_benefits"`
		Is_Intellect   int    `json:"is_intellect"`
	}

	// Parse JSON body
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

	// Kiểm tra thông tin người dùng trong token
	if claims.Uuid == "" || claims.UserName == "" || claims.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu thông tin người dùng trong token"})
		return
	}

	// Gọi gRPC đến IdeaService
	grpcConn := grpc.GetInstance()
	client := pbIdeas.NewIdeaServiceClient(grpcConn.IdeasConnect)

	res, err := client.CreateIdea(context.Background(), &pbIdeas.CreateIdeaRequest{
		IdeasName:     req.IdeasName,
		Industry:      req.Industry,
		ContentDetail: req.ContentDetail,
		Price:         req.Price,
		CustomerUuid:  claims.Uuid,
		CustomerName:  claims.UserName,
		CustomerEmail: claims.Email,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Trả về kết quả
	ctx.JSON(http.StatusOK, gin.H{
		"message": res.Message,
		"uuid":    res.Uuid,
	})
}

// //////////////////////////////////////////////////////////////////////////////
func (userCtl UserController) MyProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClaims, ok := claims.(*util.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userModel := models.Customer{}
	cond := bson.M{"uuid": userClaims.Uuid}

	user, err := userModel.FindOne(cond)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":      user.Uuid,
		"username":  user.UserName,
		"email":     user.Email,
		"image":     user.Image,
		"introduce": user.Introduce,
		"is_active": user.IsActive,
		"start_day": user.StartDay,
		"end_day":   user.EndDay,
	})
}

func (userCtl UserController) CreateRating(c *gin.Context) {
	var input struct {
		ExpertUuid string `json:"expert_uuid" binding:"required"`
		Rating     int    `json:"rating" binding:"required,min=1,max=5"`
		Comment    string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Lấy customer_id từ JWT
	customerUuid, exists := c.Get("customer_uuid")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid is missing"})
		return
	}

	customerUuidStr, ok := customerUuid.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}
	customerName, exists := c.Get("customer_name")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_name is missing"})
		return
	}
	customerNameStr, ok := customerName.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_uuid must be string"})
		return
	}
	rating := models.Rating{
		Uuid:         util.GenerateUUID(),
		CustomerUuid: customerUuidStr,
		CustomerName: customerNameStr,
		ExpertUuid:   input.ExpertUuid,
		Rating:       input.Rating,
		Comment:      input.Comment,
		CreatedAt:    time.Now(),
		IsDelete:     0,
	}

	_, err := rating.Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu đánh giá"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đánh giá đã được ghi nhận"})
}

// //////////////////////////////////////////////////////////////////////////
func (userCtl UserController) ListRating(c *gin.Context) {
	userModel := new(models.Rating)
	var req request.GetListRatingRequest
	var reqUri request.ExpertUriParam

	// BIND PATH PARAM
	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// KIỂM TRA expert_uuid tồn tại
	existCondition := bson.M{"expert_uuid": reqUri.ExpertUuid}
	_, err := userModel.FindOne(existCondition)
	if err != nil {
		c.JSON(http.StatusNotFound, respond.ErrorCommon("Không tìm thấy đánh giá nào cho chuyên gia này"))
		return
	}

	// BIND QUERY PARAM
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.MissingParams())
		return
	}

	// BUILD condition truy vấn đánh giá
	cond := bson.M{"expert_uuid": reqUri.ExpertUuid}
	if req.Rating != nil {
		cond["rating"] = *req.Rating
	}

	// PAGINATION
	optionsQuery, page, limit := models.GetPagingOption(req.Page, req.Limit, req.Sort)

	// TRUY VẤN DỮ LIỆU
	ratings, err := userModel.Pagination(c, cond, optionsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorCommon("Không thể lấy danh sách đánh giá"))
		return
	}

	// FORMAT KẾT QUẢ
	var respData []request.ListRatingResponse
	for _, rating := range ratings {
		res := request.ListRatingResponse{
			Uuid:         rating.Uuid,
			Rating:       rating.Rating,
			Comment:      rating.Comment,
			CustomerName: rating.CustomerName,
			CustomerUuid: rating.CustomerUuid,
		}
		respData = append(respData, res)
	}

	// TÍNH TỔNG
	total, err := userModel.Count(c, cond)
	if err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorCommon("Không thể đếm số lượng đánh giá"))
		return
	}
	pages := int(math.Ceil(float64(total) / float64(limit)))

	// TRẢ KẾT QUẢ
	c.JSON(http.StatusOK, respond.SuccessPagination(respData, page, limit, pages, total))
}
