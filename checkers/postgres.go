package checkers

import (
	"context"
	"errors"
	"fmt"

	"github.com/efureev/health-checker"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PostgresDBConfig struct {
	Database string
	Host     string
	Username string
	Password string
	Port     int
}

func parseDbUrl(config PostgresDBConfig) string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	return url
}

func conn(config PostgresDBConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), parseDbUrl(config))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to connect to database %s: %v\n", config.Database, err))
	}

	return conn, nil
}

func createDatabase(config PostgresDBConfig, result *checker.Result) error {
	needDb := config.Database
	config.Database = `postgres`
	step := result.AddStep(fmt.Sprintf(`Try to create database '%s'`, needDb))

	conn, err := conn(config)
	if err != nil {
		step.Failed(err)
		return err
	}

	defer conn.Close(context.Background())

	sql := fmt.Sprintf(`create database "%s" owner "%s";`, needDb, config.Username)
	_, err = conn.Exec(context.Background(), sql)
	if err != nil {
		switch err.(type) {
		case *pgconn.PgError:
			pgErr := fmt.Errorf("User \"%s\" has no permissions to create a database.\n\t%s", config.Username, err.(*pgconn.PgError))
			step.Failed(pgErr)
			return pgErr
		default:
			step.Failed(err)
			return err
		}
	}

	step.Success()

	return nil
}

func isExistDatabase(config PostgresDBConfig, result *checker.Result) error {
	step := result.AddStep(fmt.Sprintf(`Connect to database table [%s]`, config.Database))

	conn, err := conn(config)
	if err != nil {
		step.Failed(err)

		return err
	}

	conn.Close(context.Background())

	step.Success()

	return nil
}

func CheckingPostgresDatabaseFn(config PostgresDBConfig) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		err := isExistDatabase(config, result)
		if err != nil {
			needDb := config.Database
			config.Database = `postgres`
			errPg := isExistDatabase(config, result)
			if errPg != nil {
				result.AddError(fmt.Errorf(`you should create database '%s' manually`, needDb))
				return
			}

			config.Database = needDb

			errCreatedDb := createDatabase(config, result)
			if errCreatedDb != nil {
				result.AddError(fmt.Errorf("You should create database '%s' manually. \n\t%s", needDb, errCreatedDb))
				return
			}
		}

		conn, err := conn(config)
		defer conn.Close(context.Background())
		_ = conn.QueryRow(context.Background(), "SELECT version();").Scan(&result.Info.Version)

		return
	}
}
