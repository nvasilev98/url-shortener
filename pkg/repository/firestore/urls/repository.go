package urls

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
	firestoreClient *firestore.Client
}

// NewRepository is a constructor function
func NewRepository(firestoreClient *firestore.Client) *Repository {
	return &Repository{
		firestoreClient: firestoreClient,
	}
}

// AddURLTx creates URL document in firestore, if
func (r *Repository) AddURLTx(tx *firestore.Transaction, id string, url URL) error {
	doc := r.urlsCollection().Doc(id)
	err := tx.Create(doc, url)
	if err != nil {
		return fmt.Errorf("failed to create shortened url: %w", err)
	}

	return nil
}

// GetByShortURL returns a URL document by short url
// If it does not exist, it returns not found error
func (r *Repository) GetByShortURL(ctx context.Context, shortURL string) (URL, error) {
	doc, err := r.urlsCollection().Doc(shortURL).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return URL{}, NewNotFoundError()
		}

		return URL{}, fmt.Errorf("failed to retrieve by short url: %w", err)
	}

	var url URL
	if err := doc.DataTo(&url); err != nil {
		return URL{}, fmt.Errorf("failed to convert url: %w", err)
	}

	return url, nil
}

// GetDocIDByLongURL returns a URL document id by long url
// If it does not exist, it returns not found error
func (r *Repository) GetDocIDByLongURL(ctx context.Context, longURL string) (string, error) {
	collection := r.urlsCollection().
		Where("long_url", "==", longURL).
		Documents(ctx)
	defer collection.Stop()

	doc, err := collection.Next()
	if err != nil {
		if err == iterator.Done || status.Code(err) == codes.NotFound {
			return "", NewNotFoundError()
		}

		return "", fmt.Errorf("failed to retrieve by long url: %w", err)
	}

	return doc.Ref.ID, nil
}

// RunTransaction the function in a transaction
func (r *Repository) RunTransaction(ctx context.Context, txFunc func(context.Context, *firestore.Transaction) error) error {
	return r.firestoreClient.RunTransaction(ctx, txFunc)
}

func (r *Repository) urlsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("urls")
}
