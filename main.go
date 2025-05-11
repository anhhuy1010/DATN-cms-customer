package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/anhhuy1010/DATN-cms-customer/database"
	grpcClient "github.com/anhhuy1010/DATN-cms-customer/grpc"

	"github.com/anhhuy1010/DATN-cms-customer/grpc/service"

	"github.com/anhhuy1010/DATN-cms-customer/routes"
	"github.com/anhhuy1010/DATN-cms-customer/services/logService"

	pbUserr "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/users"
	"github.com/gin-gonic/gin"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	engine       *gin.Engine
	cfg          *viper.Viper
	logrusLogger *logrus.Logger
	customFunc   grpc_logrus.CodeToLevel
)

func init() {
	engine = gin.New()
	engine.Use(gin.Logger())

	logService.NewLogrus()

	cfg = config.GetConfig()
}

func main() {
	_, err := database.Init()
	if err == nil {
		fmt.Println("\n Database connected!")
	} else {
		fmt.Println("Fatal error database connection", err)
	}

	grpcSV := grpcClient.GrpcService{}
	_, err2 := grpcSV.NewService()
	if err2 == nil {
		fmt.Println("starting HTTP/2 gRPC server")
		fmt.Println()
	} else {
		fmt.Println("Fatal error GRPC connection: ", err2)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback cho local dev
	}
	go func() {
		StartRest(port)
	}()

	GRPCPort := cfg.GetString("server.grpc_port")
	err2 = StartGRPC(GRPCPort)
	if err2 != nil {
		fmt.Println("Fatal error GRPC connection: ", err2)
	}
}

func StartRest(port string) {
	routes.RouteInit(engine)

	if err := engine.Run(":" + port); err != nil {
		log.Fatalln(err)
	}
}

func StartGRPC(port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("GRPC error")
		return err
	}
	var server *grpc.Server

	// register service
	server = grpc.NewServer()
	pbUserr.RegisterUserServer(server, service.NewUserServer())

	// start gRPC server
	fmt.Println("starting gRPC server... port: ", port)

	return server.Serve(listen)
}
