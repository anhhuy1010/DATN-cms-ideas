package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/anhhuy1010/DATN-cms-ideas/config"
	"github.com/anhhuy1010/DATN-cms-ideas/database"
	grpcClient "github.com/anhhuy1010/DATN-cms-ideas/grpc"
	"github.com/anhhuy1010/DATN-cms-ideas/helpers/util"

	"github.com/anhhuy1010/DATN-cms-ideas/routes"
	"github.com/anhhuy1010/DATN-cms-ideas/services/logService"
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

	// In ra biến môi trường BASE_API_URL để kiểm tra
	fmt.Println("BASE_API_URL =", os.Getenv("BASE_API_URL"))
	// Hoặc lấy từ config nếu bạn load qua viper (nếu config có map biến này)
	fmt.Println("Config BASE_API_URL =", cfg.GetString("BASE_API_URL"))
}

func main() {

	err := util.InitMinio()
	if err != nil {
		log.Fatalln("Failed to initialize MinIO client:", err)
	}

	_, err = database.Init()
	if err == nil {
		fmt.Println("\nDatabase connected!")
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

	// Đọc PORT từ biến môi trường (Heroku hoặc local)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback cho local dev
	}

	go func() {
		StartRest(port)
	}()

	// Đọc port grpc từ config file
	GRPCPort := cfg.GetString("server.grpc_port")
	if GRPCPort == "" {
		GRPCPort = "50051" // port grpc mặc định nếu config không có
	}

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
		fmt.Println("GRPC error: ", err)
		return err
	}

	// Khởi tạo grpc server
	server := grpc.NewServer()
	// Start grpc server
	fmt.Println("starting gRPC server... port: ", port)

	return server.Serve(listen)
}
