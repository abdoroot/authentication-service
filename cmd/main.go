package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abdoroot/authentication-service/internal/auth"
	"github.com/abdoroot/authentication-service/internal/store"
	"github.com/abdoroot/authentication-service/internal/transport"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	srv            store.Storer
	host           string = os.Getenv("DB_HOST")
	port           string = os.Getenv("DB_PORT")
	databaseName   string = os.Getenv("DB_DATABASE")
	dbUsername     string = os.Getenv("DB_USERNAME")
	dbPassword     string = os.Getenv("DB_PASSWORD")
	httpListenAddr string = ":3000"
	grpcListenAddr string = ":3001"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, dbUsername, dbPassword, databaseName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal("db: ", err)
	}

	pqStore := store.NewUserStore(db)
	srv := auth.NewAuth(pqStore)

	{
		//http transport
		ht := transport.NewHttpTransport(srv, httpListenAddr)
		go ht.Strart()
	}

	{
		//grpc transport
		gt := transport.NewGRPCTransport(srv, grpcListenAddr)
		gt.Strart()
	}
}
