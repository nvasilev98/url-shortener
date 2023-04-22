package counter_test

import (
	"context"
	"url-shortener/pkg/repository/firestore/counter"
	"url-shortener/test/fixture"

	"cloud.google.com/go/firestore"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Counter Repository", func() {
	const (
		shardNumber      = 1
		shardID          = "0"
		shardsCollection = "shards"
	)

	var (
		ctx              context.Context
		firestoreClient  *firestore.Client
		repository       *counter.Repository
		firestoreFixture *fixture.FirestoreFixture
		err              error
	)

	BeforeEach(func() {
		ctx = context.Background()
		firestoreClient, err = firestore.NewClient(ctx, firestore.DetectProjectID)
		Expect(err).NotTo(HaveOccurred())
		repository = counter.NewRepository(firestoreClient, shardNumber)
		firestoreFixture = fixture.NewFirestoreFixture(firestoreClient)
	})

	AfterEach(func() {
		firestoreClient.Close()
	})

	When("it fails to initialize counter", func() {
		BeforeEach(func() {
			firestoreClient.Close()
		})

		It("should return an error", func() {
			err = repository.InitCounter(ctx)
			Expect(err).To(HaveOccurred())
		})
	})

	When("initialize counter", func() {
		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, shardsCollection, shardID)).To(Succeed())
		})

		It("should succeed", func() {
			err = repository.InitCounter(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("incrementing a shard", func() {
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, shardsCollection, shardID, counter.Shard{0})).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, shardsCollection, shardID)).To(Succeed())
		})

		It("should succeed", func() {
			err = firestoreFixture.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				return repository.IncrementCounterTx(tx)
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("it fails to convert document when getting total count", func() {
		BeforeEach(func() {
			invalidShard := map[string]interface{}{"count": "invalid"}
			Expect(firestoreFixture.InsertDocument(ctx, shardsCollection, shardID, invalidShard)).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, shardsCollection, shardID)).To(Succeed())
		})

		It("should return an error", func() {
			err = firestoreFixture.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				_, err := repository.GetCountTx(tx)
				return err
			})
			Expect(err).To(HaveOccurred())
		})
	})

	When("getting total count succeed", func() {
		expectedCount := 2
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, shardsCollection, shardID, counter.Shard{int64(expectedCount)})).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, shardsCollection, shardID)).To(Succeed())
		})

		It("should return total count", func() {
			err = firestoreFixture.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				total, err := repository.GetCountTx(tx)
				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(Equal(int64(expectedCount)))
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
