package urlshortener_test

import (
	"context"
	"errors"
	"url-shortener/cmd/urlshortener/internal/urlshortener"
	"url-shortener/cmd/urlshortener/internal/urlshortener/mocks"
	"url-shortener/pkg/repository/firestore/urls"

	"cloud.google.com/go/firestore"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Controller", func() {
	const (
		longURL  = "long-url"
		shortURL = "short-url"
	)

	var (
		mockCtrl       *gomock.Controller
		mockRepository *mocks.MockRepository
		mockCounter    *mocks.MockCounter
		mockEncoder    *mocks.MockEncoder
		controller     *urlshortener.URLController
		ctx            context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockRepository(mockCtrl)
		mockCounter = mocks.NewMockCounter(mockCtrl)
		mockEncoder = mocks.NewMockEncoder(mockCtrl)
		controller = urlshortener.NewController(mockRepository, mockCounter, mockEncoder)
		ctx = context.Background()
	})

	When("when getting document id by long url fails", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return("", errors.New("err"))
		})

		It("should return an error", func() {
			_, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).To(HaveOccurred())
		})
	})

	When("getting document id by long url succeeds", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return(shortURL, nil)
		})

		It("should return a short url", func() {
			url, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(url).To(Equal(shortURL))
		})
	})

	When("getting total count fails", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return(shortURL, urls.NewNotFoundError())
			mockRepository.EXPECT().RunTransaction(ctx, gomock.Any()).DoAndReturn(triggerTransaction)
			mockCounter.EXPECT().GetCountTx(gomock.Any()).Return(int64(0), errors.New("err"))
		})

		It("should return an error", func() {
			_, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).To(HaveOccurred())
		})
	})

	When("incrementing counter fails", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return(shortURL, urls.NewNotFoundError())
			mockRepository.EXPECT().RunTransaction(ctx, gomock.Any()).DoAndReturn(triggerTransaction)
			mockCounter.EXPECT().GetCountTx(gomock.Any()).Return(int64(0), nil)
			mockCounter.EXPECT().IncrementCounterTx(gomock.Any()).Return(errors.New("err"))
		})

		It("should return an error", func() {
			_, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).To(HaveOccurred())
		})
	})

	When("adding url fails", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return(shortURL, urls.NewNotFoundError())
			mockRepository.EXPECT().RunTransaction(ctx, gomock.Any()).DoAndReturn(triggerTransaction)
			mockCounter.EXPECT().GetCountTx(gomock.Any()).Return(int64(0), nil)
			mockCounter.EXPECT().IncrementCounterTx(gomock.Any()).Return(nil)
			mockEncoder.EXPECT().EncodeToBase62(gomock.Any()).Return(shortURL)
			mockRepository.EXPECT().AddURLTx(gomock.Any(), shortURL, urls.URL{LongURL: longURL}).Return(errors.New("err"))
		})

		It("should return an error", func() {
			_, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).To(HaveOccurred())
		})
	})

	When("running transaction succeeds", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetDocIDByLongURL(ctx, longURL).Return(shortURL, urls.NewNotFoundError())
			mockRepository.EXPECT().RunTransaction(ctx, gomock.Any()).DoAndReturn(triggerTransaction)
			mockCounter.EXPECT().GetCountTx(gomock.Any()).Return(int64(0), nil)
			mockCounter.EXPECT().IncrementCounterTx(gomock.Any()).Return(nil)
			mockEncoder.EXPECT().EncodeToBase62(gomock.Any()).Return(shortURL)
			mockRepository.EXPECT().AddURLTx(gomock.Any(), shortURL, urls.URL{LongURL: longURL}).Return(nil)
		})

		It("should return short url", func() {
			url, err := controller.CreateShortURL(ctx, longURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(url).To(Equal(shortURL))
		})
	})

	When("getting url object by short url succeds", func() {
		BeforeEach(func() {
			mockRepository.EXPECT().GetByShortURL(ctx, shortURL).Return(urls.URL{LongURL: longURL}, nil)
		})

		It("should return url object with long address", func() {
			url, err := controller.GetByShortURL(ctx, shortURL)
			Expect(err).ToNot(HaveOccurred())
			Expect(url.LongURL).To(Equal(longURL))
		})
	})

})

func triggerTransaction(ctx context.Context, txFunc func(context.Context, *firestore.Transaction) error) error {
	return txFunc(ctx, &firestore.Transaction{})
}
