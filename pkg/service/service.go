package service

import (
	"URLShortener/pkg/repository"
)

type UrlList interface {
	GetUrl(shortUrl string) (string, error)
	PostUrl(longUrl string) (string, error)
}

type Service struct {
	UrlList
}

func NewService(repo *repository.Repository, stringLength int, characters []rune) *Service {
	return &Service{
		UrlList: NewUrlListService(repo.UrlList, stringLength, characters),
	}
}
