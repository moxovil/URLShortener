package repository

// var ErrNoRecord = errors.New("подходящей записи не найдено")
// var ErrExistRecord = errors.New("запись уже существует")

type UrlList interface {
	GetUrl(shortUrl string) (string, error)
	GetShortUrl(longUrl string) (string, error)
	PostUrl(shortUrl string, longUrl string) error
	IsAvailable(shortUrl string) (bool, error)
}

type Repository struct {
	UrlList
}

func NewRepository(db UrlList) *Repository {
	return &Repository{
		UrlList: db,
	}
}
