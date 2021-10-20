package inmemory

import (
	"errors"
)

type InMemoryDB struct {
	longUrlsMap map[string]string
}

func NewInMemoryDB() *InMemoryDB {
	longUrlsMap := make(map[string]string)
	return &InMemoryDB{
		longUrlsMap: longUrlsMap,
	}
}

func (db *InMemoryDB) GetUrl(shortUrl string) (string, error) {
	if shortUrl == "" {
		return "", errors.New("неправильные исходные данные для чтения из БД")
	}
	longUrl, ok := db.longUrlsMap[shortUrl]
	if ok {
		return longUrl, nil
	}
	return "", nil
}

func (db *InMemoryDB) GetShortUrl(longUrl string) (string, error) {
	if longUrl == "" {
		return "", errors.New("неправильные исходные данные для чтения из БД")
	}
	for key, value := range db.longUrlsMap {
		if longUrl == value {
			return key, nil
		}
	}
	return "", nil
}

func (db *InMemoryDB) PostUrl(shortUrl string, longUrl string) error {
	if shortUrl == "" || longUrl == "" {
		return errors.New("неправильные исходные данные для записи в БД")
	}

	isExistShortUrl, _ := db.GetShortUrl(longUrl)
	if isExistShortUrl != "" {
		return errors.New("данные уже записаны в БД")
	}

	isExistLongUrl, _ := db.GetUrl(shortUrl)
	if isExistLongUrl != "" {
		return errors.New("данные уже записаны в БД")
	}

	db.longUrlsMap[shortUrl] = longUrl
	return nil
}

func (db *InMemoryDB) IsAvailable(shortUrl string) (bool, error) {
	if shortUrl == "" {
		return false, errors.New("неправильные исходные данные для записи в БД")
	}
	_, ok := db.longUrlsMap[shortUrl]
	return !ok, nil
}
