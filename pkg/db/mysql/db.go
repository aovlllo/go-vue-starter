package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
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

// NewConnection creates a new database connection
func NewConnection(config *Config) (*DB, error) {
	conn := config.ConnectionURL()
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

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
