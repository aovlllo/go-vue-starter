package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

// Config represents MySQL configuration
type Config struct {
	Name     string `yaml:"db_name"`
	Host     string `yaml:"db_host"`
	Port     string `yaml:"db_port"`
	User     string `yaml:",omitempty" env:"DB_USER"`
	Password string `yaml:",omitempty" env:"DB_PASSWORD"`
}

// DB represents the structure of the database
type DB struct {
	client *sql.DB
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if v := c.Port; v != "" {
		host = host + ":" + v
	}
	host = "tcp(" + host + ")"

	return fmt.Sprintf("%s@%s/%s", url.UserPassword(c.User, c.Password), host, c.Name)
}

func tryConnection(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewConnection creates a new database connection
func NewConnection(config *Config) (*DB, error) {
	conn := config.ConnectionURL()

	var db *sql.DB
	var err error
	for i := 1; i < 6; i++ {
		db, err = tryConnection(conn)
		if err == nil {
			break
		}
		log.Warn().Msgf("unsuccessful connection, wait for %d seconds", i*10)
		time.Sleep(time.Duration(i*10) * time.Second)
	}
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(4 * time.Minute)

	return &DB{db}, nil
}

// CloseConnection closes the database connection
func (db *DB) CloseConnection() error {
	err := db.client.Close()
	if err != nil {
		return err
	}

	return nil
}
