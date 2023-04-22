package urlshortener_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"url-shortener/cmd/urlshortener/internal/urlshortener"
	"url-shortener/cmd/urlshortener/internal/urlshortener/mocks"
	"url-shortener/pkg/repository/firestore/urls"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Presenter", func() {
	const (
		longURL  = "long-url"
		shortURL = "short-url"
	)

	var (
		mockCtrl       *gomock.Controller
		mockContext    *gin.Context
		recorder       *httptest.ResponseRecorder
		mockController *mocks.MockController
		presenter      *urlshortener.Presenter
		err            error
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		mockContext, _ = gin.CreateTestContext(recorder)
		mockCtrl = gomock.NewController(GinkgoT())
		mockController = mocks.NewMockController(mockCtrl)
		presenter = urlshortener.NewPresenter(mockController)
	})

	When("it fails to create short url", func() {
		BeforeEach(func() {
			mockContext.Request, err = http.NewRequest(http.MethodPost, gomock.Any().String(), bytes.NewBufferString(longURL))
			Expect(err).ToNot(HaveOccurred())
			mockController.EXPECT().CreateShortURL(gomock.Any(), longURL).Return("", errors.New("err"))
		})

		It("should return http status internal server error", func() {
			presenter.CreateShortURL(mockContext)
			Expect(mockContext.Writer.Status()).To(Equal(http.StatusInternalServerError))
		})
	})

	When("it succeeds to create a short url", func() {
		BeforeEach(func() {
			mockContext.Request, err = http.NewRequest(http.MethodPost, gomock.Any().String(), bytes.NewBufferString(longURL))
			Expect(err).ToNot(HaveOccurred())
			mockController.EXPECT().CreateShortURL(gomock.Any(), longURL).Return(shortURL, nil)
		})

		It("should return http status ok and short url", func() {
			presenter.CreateShortURL(mockContext)
			Expect(mockContext.Writer.Status()).To(Equal(http.StatusOK))
			Expect(recorder.Body.String()).To(ContainSubstring(shortURL))
		})
	})

	When("it fails to get by short url", func() {
		BeforeEach(func() {
			mockContext.Request, err = http.NewRequest(http.MethodGet, gomock.Any().String(), nil)
			Expect(err).ToNot(HaveOccurred())
			mockContext.Params = []gin.Param{{Key: "short_url", Value: shortURL}}
			mockController.EXPECT().GetByShortURL(gomock.Any(), shortURL).Return(urls.URL{}, errors.New("err"))
		})

		It("should return http status internal server error", func() {
			presenter.RedirectToLongURL(mockContext)
			Expect(mockContext.Writer.Status()).To(Equal(http.StatusInternalServerError))
		})
	})

	When("short url does not exist", func() {
		BeforeEach(func() {
			mockContext.Request, err = http.NewRequest(http.MethodGet, gomock.Any().String(), nil)
			Expect(err).ToNot(HaveOccurred())
			mockContext.Params = []gin.Param{{Key: "short_url", Value: shortURL}}
			mockController.EXPECT().GetByShortURL(gomock.Any(), shortURL).Return(urls.URL{}, urls.NewNotFoundError())
		})

		It("should return http status not found", func() {
			presenter.RedirectToLongURL(mockContext)
			Expect(mockContext.Writer.Status()).To(Equal(http.StatusNotFound))
		})
	})

	When("short url is found", func() {
		BeforeEach(func() {
			mockContext.Request, err = http.NewRequest(http.MethodGet, gomock.Any().String(), nil)
			Expect(err).ToNot(HaveOccurred())
			mockContext.Params = []gin.Param{{Key: "short_url", Value: shortURL}}
			mockController.EXPECT().GetByShortURL(gomock.Any(), shortURL).Return(urls.URL{LongURL: longURL}, nil)
		})

		It("should return status found and redirect to long url", func() {
			presenter.RedirectToLongURL(mockContext)
			Expect(mockContext.Writer.Status()).To(Equal(http.StatusFound))
			Expect(recorder.Body.String()).To(ContainSubstring(longURL))
		})
	})
})
