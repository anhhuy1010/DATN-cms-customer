package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/anhhuy1010/DATN-cms-customer/config"
)

var db *mongo.Database

func Init() (*mongo.Database, error) {
	if db == nil {
		cfg := config.GetConfig()
		host := cfg.GetString("database.host")
		port := cfg.GetString("database.port")
		user := cfg.GetString("database.username")
		password := cfg.GetString("database.password")
		database := cfg.GetString("database.db_name")
		ssl := cfg.GetBool("database.ssl")

		var uri string
		if ssl {
			// URI MongoDB Atlas (mongodb+srv)
			uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&readPreference=secondaryPreferred",
				user, password, host)
		} else {
			// URI MongoDB local hoặc không dùng TLS
			uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?authMechanism=SCRAM-SHA-256",
				user, password, host, port)
		}

		fmt.Printf("🔗 MongoDB URI: %s\n", uri) // debug URI

		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()
		clientOptions := options.Client().ApplyURI(uri)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			fmt.Println("❌ Kết nối MongoDB thất bại:", err)
			return nil, err
		}

		if err := client.Ping(ctx, nil); err != nil {
			fmt.Println("❌ Ping MongoDB thất bại:", err)
			return nil, err
		}

		db = client.Database(database)
		fmt.Println("✅ Đã kết nối MongoDB thành công.")
	}

	return db, nil
}

func GetInstance() *mongo.Database {
	return db
}
