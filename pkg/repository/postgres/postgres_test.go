package postgres

import (
	"fmt"
	"testing"
)

// host: "db"
// port: "5432"
// user: "postgres"
// password: "qwerty"
// dbname: "postgres"
// sslmode: "disable"

func TestPostUrl(t *testing.T) {
	db, err := NewPostgresDB(Config{
		Host:     "db",
		Port:     "5432",
		User:     "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	if err != nil {
		t.Error(fmt.Sprintf("ошибка инициализации БД: %s\n", err.Error()))
	}

	defer db.Close()

	type postTestData struct {
		shortUrl string
		longUrl  string
		err      bool
	}

	var postTests = []postTestData{
		{"val1", "val2", false},
		{"", "val3", true},
		{"val4", "", true},
		{"val5", "val2", true},
	}

	for _, postTest := range postTests {
		err := db.PostUrl(postTest.shortUrl, postTest.longUrl)
		if (err != nil) != postTest.err {
			t.Error(fmt.Sprintf("PostUrl(\"%s\", \"%s\") = %t, expected %t", postTest.shortUrl, postTest.longUrl, err != nil, postTest.err))
		}
	}
}

func TestGetUrl(t *testing.T) {
	db, err := NewPostgresDB(Config{
		Host:     "db",
		Port:     "5432",
		User:     "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	if err != nil {
		t.Error(fmt.Sprintf("ошибка инициализации БД: %s\n", err.Error()))
	}

	defer db.Close()

	type initialData struct {
		shortUrl string
		longUrl  string
	}

	var initialDataToPost = []initialData{
		{"var1", "var2"},
		{"var3", "var4"},
	}

	type getTestData struct {
		shortUrl string
		longUrl  string
		err      bool
	}

	var getTests = []getTestData{
		{"var1", "var2", false},
		{"", "", true},
		{"var3", "var4", false},
		{"var5", "", false},
	}

	for _, dataToPost := range initialDataToPost {
		db.PostUrl(dataToPost.shortUrl, dataToPost.longUrl)
	}

	for _, getTest := range getTests {
		longUrl, err := db.GetUrl(getTest.shortUrl)
		if (longUrl != getTest.longUrl) || (err != nil) != getTest.err {
			t.Error(fmt.Sprintf("GetUrl(\"%s\") = \"%s\", %t expected: \"%s\", %t", getTest.shortUrl, longUrl, err != nil, getTest.longUrl, getTest.err))
		}
	}
}

func TestGetShortUrl(t *testing.T) {
	db, err := NewPostgresDB(Config{
		Host:     "db",
		Port:     "5432",
		User:     "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	if err != nil {
		t.Error(fmt.Sprintf("ошибка инициализации БД: %s\n", err.Error()))
	}

	defer db.Close()

	type initialData struct {
		shortUrl string
		longUrl  string
	}

	var initialDataToPost = []initialData{
		{"var1", "var2"},
		{"var3", "var4"},
	}

	type getShortTestData struct {
		shortUrl string
		longUrl  string
		err      bool
	}

	var getShortTests = []getShortTestData{
		{"var1", "var2", false},
		{"", "", true},
		{"var3", "var4", false},
		{"", "var5", false},
	}

	for _, dataToPost := range initialDataToPost {
		db.PostUrl(dataToPost.shortUrl, dataToPost.longUrl)
	}

	for _, getShortTest := range getShortTests {
		shortUrl, err := db.GetShortUrl(getShortTest.longUrl)
		if (shortUrl != getShortTest.shortUrl) || (err != nil) != getShortTest.err {
			t.Error(fmt.Sprintf("GetShortUrl(\"%s\") = \"%s\", %t expected: \"%s\", %t", getShortTest.longUrl, shortUrl, err != nil, getShortTest.shortUrl, getShortTest.err))
		}
	}
}

func TestIsAvailable(t *testing.T) {
	db, err := NewPostgresDB(Config{
		Host:     "db",
		Port:     "5432",
		User:     "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})

	if err != nil {
		t.Error(fmt.Sprintf("ошибка инициализации БД: %s\n", err.Error()))
	}

	defer db.Close()

	type initialData struct {
		shortUrl string
		longUrl  string
	}

	var initialDataToPost = []initialData{
		{"var1", "var2"},
		{"var3", "var4"},
	}

	type availableTestData struct {
		shortUrl string
		ok       bool
		err      bool
	}

	var availableTests = []availableTestData{
		{"var1", false, false},
		{"", false, true},
		{"var3", false, false},
		{"var5", true, false},
	}

	for _, dataToPost := range initialDataToPost {
		db.PostUrl(dataToPost.shortUrl, dataToPost.longUrl)
	}

	for _, availableTest := range availableTests {
		ok, err := db.IsAvailable(availableTest.shortUrl)
		if (ok != availableTest.ok) || (err != nil) != availableTest.err {
			t.Error(fmt.Sprintf("IsAvailable(\"%s\") = %t, %t expected: %t, %t", availableTest.shortUrl, ok, err != nil, availableTest.ok, availableTest.err))
		}
	}
}
