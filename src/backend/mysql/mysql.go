package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

const (
	defaultPortForMysql = 3306
	defaultPingInterval = 10
	defaultWaitTimeout  = 7200
	maxIdleConnections  = 10
)

type Client struct {
	*sql.DB

	logger *zap.Logger
}

//New returns backend client with error to the calling function.
func New(logger *zap.Logger, conf *types.DatabaseMetadata, modelsDir string) (*Client, error) {
	if conf.Port == 0 {
		conf.Port = defaultPortForMysql
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?collation=utf8mb4_general_ci", conf.User, conf.Password, conf.Host, conf.Port, conf.Name))
	if err != nil {
		return nil, err
	}
	err = bootstrapDB(db, modelsDir)
	if err != nil {
		return nil, err
	}
	err = optimizeDBClientConnection(db)
	if err != nil {
		return nil, err
	}
	return &Client{DB: db, logger: logger}, err
}

func bootstrapDB(db *sql.DB, modelsDir string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", modelsDir), "mysql", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func optimizeDBClientConnection(db *sql.DB) error {
	//https://www.alexedwards.net/blog/configuring-sqldb
	//good read to see what values for optimization of the backend
	//connections be tuned to get the best results
	maxConnectionLifetime := time.Duration(maxIdleConnections)
	waitTimeout := "wait_timeout"
	stmt, err := db.Query("show variables where variable_name = ?", waitTimeout)
	defer func() {
		if err := stmt.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing db statement, err: %+v", err)
		}
	}()
	if err != nil {
		return err
	}
	if stmt.Next() {
		err := stmt.Scan(&waitTimeout, &maxConnectionLifetime)
		if err != nil {
			return err
		}
		maxConnectionLifetime = maxConnectionLifetime / 2
	}
	maxConnectionLifetime = maxConnectionLifetime * time.Second
	db.SetConnMaxLifetime(maxConnectionLifetime)
	db.SetMaxIdleConns(defaultWaitTimeout)
	return nil
}

func (c *Client) Close() {
	err := c.DB.Close()
	if err != nil {
		c.logger.Error("could not close connection to db", zap.Error(err))
	}
}
