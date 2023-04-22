package urlshortener

import (
	"context"
	"errors"
	"fmt"
	"url-shortener/pkg/repository/firestore/urls"

	"cloud.google.com/go/firestore"
)

//go:generate mockgen --source=controller.go --destination mocks/controller.go --package mocks

type Repository interface {
	AddURLTx(tx *firestore.Transaction, id string, url urls.URL) error
	GetByShortURL(ctx context.Context, shortURL string) (urls.URL, error)
	GetDocIDByLongURL(ctx context.Context, longURL string) (string, error)
	RunTransaction(ctx context.Context, txFunc func(context.Context, *firestore.Transaction) error) error
}

type Counter interface {
	IncrementCounterTx(tx *firestore.Transaction) error
	GetCountTx(tx *firestore.Transaction) (int64, error)
}

type Encoder interface {
	EncodeToBase62(number uint64) string
}

type URLController struct {
	repository Repository
	counter    Counter
	encoder    Encoder
}

// NewController is a constructor function
func NewController(repository Repository, counter Counter, encoder Encoder) *URLController {
	return &URLController{
		repository: repository,
		counter:    counter,
		encoder:    encoder,
	}
}

// CreateShortURL creates an URL object if not exists, otherwise it returns the id of the existing one
func (c *URLController) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	id, err := c.repository.GetDocIDByLongURL(ctx, longURL)
	if err != nil {
		var notFoundErr urls.NotFoundError
		if errors.As(err, &notFoundErr) {
			return c.createShortURL(ctx, longURL)
		}

		return "", fmt.Errorf("failed to get doc id by long url: %w", err)
	}

	return id, nil
}

// GetByShortURL return URL object by short URL address
func (c *URLController) GetByShortURL(ctx context.Context, shortURL string) (urls.URL, error) {
	return c.repository.GetByShortURL(ctx, shortURL)
}

func (c *URLController) createShortURL(ctx context.Context, longURL string) (string, error) {
	var id string
	err := c.repository.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		total, err := c.counter.GetCountTx(tx)
		if err != nil {
			return err
		}

		if err := c.counter.IncrementCounterTx(tx); err != nil {
			return err
		}

		id = c.encoder.EncodeToBase62(uint64(total + 1))
		return c.repository.AddURLTx(tx, id, urls.URL{LongURL: longURL})
	})

	if err != nil {
		return "", fmt.Errorf("failed to run transaction: %w", err)
	}

	return id, nil
}
