package handler

import (
	"URLShortener/pkg/service"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	services *service.Service
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewHandler(services *service.Service, errorLog *log.Logger, infoLog *log.Logger) *Handler {
	return &Handler{
		services: services,
		errorLog: errorLog,
		infoLog:  infoLog,
	}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.shortUrlsHandler)
	return mux
}

func (h *Handler) shortUrlsHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	//Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL
	if r.Method == http.MethodGet {
		// app.infolog.Println("Метод GET:")
		// app.writeRequest(r)
		shortUrl := r.URL.Query().Get("url")
		if shortUrl == "" {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		url, err := h.services.GetUrl(shortUrl)
		if err != nil {
			h.serverError(w)
			return
		}
		if url == "" {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		w.Write([]byte(url))
		return
	}

	//Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый
	if r.Method == http.MethodPost {
		// app.writeRequest(r)
		if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
			h.clientError(w, http.StatusUnsupportedMediaType)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.serverError(w)
			return
		}

		url := string(body)
		if url == "" {
			h.clientError(w, http.StatusBadRequest)
			return
		}

		shortUrl, err := h.services.PostUrl(url)

		if err != nil {
			if shortUrl != "" {
				h.clientError(w, http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Для URL: %s уже имеется короткая запись: %s\n", url, shortUrl)))
				return
			}
			h.serverError(w)
			return
		}

		w.Write([]byte(shortUrl))
		return
	}

	w.Header().Set("Allow", http.MethodPost+"; "+http.MethodGet)
	h.clientError(w, http.StatusMethodNotAllowed)
}
