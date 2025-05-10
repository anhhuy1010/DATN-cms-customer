package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/anhhuy1010/DATN-cms-customer/config"
)

var (
	db     *mongo.Database
	client *mongo.Client
)

// Init connects to MongoDB and returns the database instance.
func Init() (*mongo.Database, error) {
	if db != nil {
		return db, nil
	}

	cfg := config.GetConfig()

	// Đọc thông tin từ config
	user := url.QueryEscape(cfg.GetString("database.username"))
	password := url.QueryEscape(cfg.GetString("database.password"))
	host := cfg.GetString("database.host")
	dbName := cfg.GetString("database.db_name")
	ssl := cfg.GetBool("database.ssl")

	// Xây dựng URI MongoDB Atlas
	var uri string
	if ssl {
		// Dùng SRV URI nếu dùng MongoDB Atlas (SSL bắt buộc)
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, password, host)
	} else {
		// (Trường hợp đặc biệt nếu dùng cluster thường không có SRV)
		log.Fatal("❌ MongoDB Atlas yêu cầu SSL. Vui lòng bật 'ssl: true' trong config.")
	}

	log.Printf("🔗 Đang kết nối MongoDB Atlas tại host: %s (DB: %s)", host, dbName)

	// Thiết lập timeout cho kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tạo client với URI
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Kết nối MongoDB thất bại: %v", err)
		return nil, err
	}

	// Ping kiểm tra kết nối
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Ping MongoDB thất bại: %v", err)
		return nil, err
	}

	// Gán database và trả về
	db = client.Database(dbName)
	log.Println("✅ Kết nối MongoDB thành công.")
	return db, nil
}

// GetInstance trả về instance đang kết nối
func GetInstance() *mongo.Database {
	if db == nil {
		log.Fatal("❌ MongoDB chưa được khởi tạo. Gọi Init() trước.")
	}
	return db
}

// Disconnect đóng kết nối MongoDB
func Disconnect() {
	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Println("⚠️ Không thể ngắt kết nối MongoDB:", err)
		} else {
			log.Println("🔌 Đã ngắt kết nối MongoDB.")
		}
	}
}
