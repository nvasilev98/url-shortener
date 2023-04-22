package fixture

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type FirestoreFixture struct {
	client *firestore.Client
}

func NewFirestoreFixture(client *firestore.Client) *FirestoreFixture {
	return &FirestoreFixture{client}
}

func (f *FirestoreFixture) InsertDocument(ctx context.Context, collection, id string, data interface{}) error {
	docRef := f.client.Doc(fmt.Sprintf("%s/%s", collection, id))
	_, err := docRef.Set(ctx, data)
	return err
}

func (f *FirestoreFixture) DeleteDocument(ctx context.Context, collection, id string) error {
	docRef := f.client.Doc(fmt.Sprintf("%s/%s", collection, id))
	_, err := docRef.Delete(ctx)
	return err
}

func (f *FirestoreFixture) RunTransaction(ctx context.Context, txFunc func(ctx context.Context, tx *firestore.Transaction) error) error {
	return f.client.RunTransaction(ctx, txFunc)
}
