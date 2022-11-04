package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-sms/internal/jsonlog"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "github.com/twilio/twilio-go"
	// verify "github.com/twilio/twilio-go/rest/verify/v2"
)

const Version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintError(err, nil)
		return
	}

	logger.PrintInfo("Environment variables loaded successfully", nil)

	var cfg config

	flag.IntVar(&cfg.port, "port", 4005, "API server port")
	flag.StringVar(&cfg.env, "environment", "development", "Environment | development | staging | production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_DSN"), "database connection string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL maximum open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL maximum idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "10m", "PostgreSQL maximum idle time")
	flag.Parse()

	db, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	app := &application{
		config: cfg,
		logger: logger,
	}

	logger.PrintInfo("stating server", map[string]string{
		"addr": fmt.Sprintf(":%d", cfg.port),
		"env":  cfg.env,
	})

	if err = app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

	// client := twilio.NewRestClient()
	// params := verify.CreateVerificationParams{}
	// params.SetTo("+254736546908")
	// params.SetChannel("sms")
	// resp, err := client.VerifyV2.CreateVerification("VAb8b1b718df0b522d53e4759de792a7f4", &params)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	if resp.Status != nil {
	// 		fmt.Println(*resp.Status)
	// 		fmt.Println(*resp.SendCodeAttempts...)
	// 		fmt.Println(resp.To)
	// 	} else {
	// 		fmt.Println(resp.Status)
	// 	}
	// }

	// params := verify.CreateServiceParams{}
	// params.SetFriendlyName("Rabitechs")

	// resp, err := client.VerifyV2.CreateService(&params)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	if resp.Sid != nil {
	// 		fmt.Println(*resp.Sid)
	// 	} else {
	// 		fmt.Println(resp.Sid)
	// 	}
	// }

	// fmt.Println(*resp.Links)

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return db, nil
}
