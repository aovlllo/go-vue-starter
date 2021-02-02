package db

import (
	"github.com/aovlllo/vue-template/pkg/db/mysql"
	"github.com/aovlllo/vue-template/pkg/model"
)


// Config represents the configuration of the database interface
type Config struct {
	MySQL *mysql.Config
}

// DB is the interface which must be implemented by all db drivers
type DB interface {
	CloseConnection() error

	CreateUser(u *model.User) error
	GetUser(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(u *model.User) error
	DeleteUser(id int) error
}

// NewConnection creates a new database connection
func NewConnection(config *Config) (DB, error) {
	db, err := mysql.NewConnection(config.MySQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
