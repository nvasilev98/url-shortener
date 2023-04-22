package counter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repository struct {
	firestoreClient *firestore.Client
	ShardsNumber    int
}

// NewRepository is a constructor function
func NewRepository(firestoreClient *firestore.Client, shardsNumber int) *Repository {
	return &Repository{
		firestoreClient: firestoreClient,
		ShardsNumber:    shardsNumber,
	}
}

// InitCounter creates a given number of shards as subcollection
// It sets to each shard an initial value to 0 if it does not exist
func (r *Repository) InitCounter(ctx context.Context) error {
	collectionRef := r.shardsCollection()
	for num := 0; num < r.ShardsNumber; num++ {
		if _, err := collectionRef.Doc(strconv.Itoa(num)).Create(ctx, Shard{Count: 0}); err != nil {
			if status.Code(err) == codes.AlreadyExists {
				continue
			}

			return fmt.Errorf("failed to create shard with id [%d]: %w", num, err)
		}
	}

	return nil
}

// IncrementCounterTx increments a random shard
func (r *Repository) IncrementCounterTx(tx *firestore.Transaction) error {
	docID := strconv.Itoa(rand.Intn(r.ShardsNumber))
	shardRef := r.shardsCollection().Doc(docID)
	if err := tx.Update(shardRef, []firestore.Update{{Path: "count", Value: firestore.Increment(1)}}); err != nil {
		return fmt.Errorf("failed to update shard: %w", err)
	}

	return nil
}

// GetCountTx get total count across all shards
func (r *Repository) GetCountTx(tx *firestore.Transaction) (int64, error) {
	var total int64
	shards := tx.Documents(r.shardsCollection())
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return 0, fmt.Errorf("failed to return next shard: %w", err)
		}

		var shard Shard
		if err = doc.DataTo(&shard); err != nil {
			return 0, fmt.Errorf("failed to convert shard: %w", err)
		}

		total += shard.Count
	}

	return total, nil
}

func (r *Repository) shardsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("shards")
}
