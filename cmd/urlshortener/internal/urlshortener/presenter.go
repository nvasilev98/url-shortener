package urlshortener

import (
	"context"
	"errors"
	"net/http"
	"url-shortener/pkg/repository/firestore/urls"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen --source=presenter.go --destination mocks/presenter.go --package mocks

type Controller interface {
	CreateShortURL(ctx context.Context, longURL string) (string, error)
	GetByShortURL(ctx context.Context, shortURL string) (urls.URL, error)
}

type Presenter struct {
	controller Controller
}

// NewPresenter is a constructor function
func NewPresenter(controller Controller) *Presenter {
	return &Presenter{
		controller: controller,
	}
}

// CreateShortURL creates a short URL object and returns its ID
func (p *Presenter) CreateShortURL(ctx *gin.Context) {
	urlAddress, err := ctx.GetRawData()
	if err != nil {
		logrus.Errorf("Failed to read body %v", err)
		ctx.JSON(http.StatusInternalServerError, "Error occurred while reading body")
		return
	}

	shortID, err := p.controller.CreateShortURL(ctx, string(urlAddress))
	if err != nil {
		logrus.Errorf("Failed to create short url: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Error occured while creating short URL")
		return
	}

	ctx.JSON(http.StatusOK, shortID)
}

// RedirectToLongURL accepts a short URL as path param and redirects to the long URL if it exists
func (p *Presenter) RedirectToLongURL(ctx *gin.Context) {
	shortURL := ctx.Param("short_url")
	url, err := p.controller.GetByShortURL(ctx, shortURL)
	if err != nil {
		var notFoundErr urls.NotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.JSON(http.StatusNotFound, "URL does not exist")
			return
		}

		logrus.Errorf("Failed to get by short url: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Error occured while getting short URL")
		return
	}

	ctx.Redirect(http.StatusFound, url.LongURL)
}
