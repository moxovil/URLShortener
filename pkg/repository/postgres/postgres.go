package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

const urlsTable = "urls"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (short_url varchar(255) PRIMARY KEY, long_url varchar(255) not null unique);", urlsTable)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (postgres *PostgresDB) GetUrl(shortUrl string) (string, error) {
	if shortUrl == "" {
		return "", errors.New("неправильные исходные данные для чтения из БД")
	}
	query := fmt.Sprintf("SELECT long_url FROM %s WHERE short_url = $1", urlsTable)
	row := postgres.db.QueryRow(query, shortUrl)

	var longUrl string
	err := row.Scan(&longUrl)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return longUrl, nil
}

func (postgres *PostgresDB) GetShortUrl(longUrl string) (string, error) {
	if longUrl == "" {
		return "", errors.New("неправильные исходные данные для чтения из БД")
	}
	query := fmt.Sprintf("SELECT short_url FROM %s WHERE long_url = $1", urlsTable)
	row := postgres.db.QueryRow(query, longUrl)

	var shortUrl string
	err := row.Scan(&shortUrl)

	if err != nil {
		return "", err
	}
	return shortUrl, nil
}

func (postgres *PostgresDB) PostUrl(shortUrl string, longUrl string) error {
	if shortUrl == "" || longUrl == "" {
		return errors.New("неправильные исходные данные для записи в БД")
	}

	query := fmt.Sprintf("SELECT short_url FROM %s WHERE short_url = $1 OR long_url = $2", urlsTable)
	row := postgres.db.QueryRow(query, shortUrl, longUrl)

	var url string
	err := row.Scan(&url)

	if err == sql.ErrNoRows {
		query = fmt.Sprintf("INSERT INTO %s (short_url, long_url) values ($1, $2);", urlsTable)
		_, err = postgres.db.Query(query, shortUrl, longUrl)
		return err
	} else if err != nil {
		return err
	}

	return errors.New("данные уже записаны в БД")
}

func (postgres *PostgresDB) IsAvailable(shortUrl string) (bool, error) {
	if shortUrl == "" {
		return false, errors.New("неправильные исходные данные для записи в БД")
	}

	url, err := postgres.GetUrl(shortUrl)

	if err != nil {
		return false, err
	}

	return url == "", nil
}

func (postgres *PostgresDB) Close() error {
	return postgres.db.Close()
}
