package urls_test

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	. "github.com/onsi/ginkgo/v2"

	"url-shortener/pkg/repository/firestore/urls"
	"url-shortener/test/fixture"

	. "github.com/onsi/gomega"
)

var _ = Describe("URLs Repository", func() {
	const (
		id             = "test-id"
		longURL        = "url"
		urlsCollection = "urls"
	)

	var (
		ctx              context.Context
		firestoreClient  *firestore.Client
		repository       *urls.Repository
		firestoreFixture *fixture.FirestoreFixture
		err              error
	)

	BeforeEach(func() {
		ctx = context.Background()
		firestoreClient, err = firestore.NewClient(ctx, firestore.DetectProjectID)
		Expect(err).NotTo(HaveOccurred())
		repository = urls.NewRepository(firestoreClient)
		firestoreFixture = fixture.NewFirestoreFixture(firestoreClient)
	})

	AfterEach(func() {
		firestoreClient.Close()
	})

	When("it fails to add an url document", func() {
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, urlsCollection, id, urls.URL{})).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, urlsCollection, id)).To(Succeed())
		})

		It("should return an error", func() {
			err = firestoreFixture.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				return repository.AddURLTx(tx, id, urls.URL{})
			})

			Expect(err).To(HaveOccurred())
		})
	})

	When("creating an url object", func() {
		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, urlsCollection, id)).To(Succeed())
		})

		It("should succeed", func() {
			err = firestoreFixture.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				return repository.AddURLTx(tx, id, urls.URL{})
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("getting document fails", func() {
		BeforeEach(func() {
			firestoreClient.Close()
		})

		It("should return an error", func() {
			_, err := repository.GetByShortURL(ctx, id)
			Expect(err).To(HaveOccurred())
		})
	})

	When("getting document by id that does not exist", func() {
		It("should return an error", func() {
			_, err := repository.GetByShortURL(ctx, "unknown-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(urls.NotFoundError{}))
		})
	})

	When("it fails to convert document to url object", func() {
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, urlsCollection, id, map[string]interface{}{"long_url": time.Now()})).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, urlsCollection, id)).To(Succeed())
		})

		It("should return an error", func() {
			_, err := repository.GetByShortURL(ctx, id)
			Expect(err).To(HaveOccurred())
		})
	})

	When("getting document by id that exists", func() {
		var expectedURL = urls.URL{LongURL: longURL}
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, urlsCollection, id, expectedURL)).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, urlsCollection, id)).To(Succeed())
		})

		It("should return an url object", func() {
			url, err := repository.GetByShortURL(ctx, id)
			Expect(err).ToNot(HaveOccurred())
			Expect(url).To(Equal(expectedURL))
		})
	})

	When("getting document id by long url fails", func() {
		BeforeEach(func() {
			firestoreClient.Close()
		})

		It("should return an error", func() {
			_, err := repository.GetDocIDByLongURL(ctx, longURL)
			Expect(err).To(HaveOccurred())
		})
	})

	When("getting document id by long url that does not exists", func() {
		It("should return an error", func() {
			_, err := repository.GetDocIDByLongURL(ctx, "unknown-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(urls.NotFoundError{}))
		})
	})

	When("getting document id by long url succeeds", func() {
		BeforeEach(func() {
			Expect(firestoreFixture.InsertDocument(ctx, urlsCollection, id, urls.URL{LongURL: longURL})).To(Succeed())
		})

		AfterEach(func() {
			Expect(firestoreFixture.DeleteDocument(ctx, urlsCollection, id)).To(Succeed())
		})

		It("should return an document id", func() {
			docID, err := repository.GetDocIDByLongURL(ctx, longURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(docID).To(Equal(id))
		})
	})
})
