package dbconfig

import (
	"context"
	"example.com/test/pkg/config"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Executor interface {
	GetConnection() (*pgxpool.Pool, error)
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SslMode  string
}

func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     config.Configuration.Database.Host,
		Port:     config.Configuration.Database.Port,
		User:     config.Configuration.Database.User,
		Password: config.Configuration.Database.Password,
		DbName:   config.Configuration.Database.DbName,
		SslMode:  config.Configuration.Database.SslMode,
	}
}

func (c *DatabaseConfig) GetConnection() (*pgxpool.Pool, error) {
	pgConn, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DbName))
	if err != nil {
		return nil, err
	}
	return pgConn, nil
}
