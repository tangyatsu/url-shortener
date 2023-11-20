package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			urlGetterMock := mocks.NewURLGetter(t)

			if test.respError == "" || test.mockError != nil {
				urlGetterMock.On("GetURL", test.alias).
					Return(test.url, test.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			tsrv := httptest.NewServer(r)
			defer tsrv.Close()

			redirectURL, err := api.GetRedirect(tsrv.URL + "/" + test.alias)
			require.NoError(t, err)

			require.Equal(t, redirectURL, test.url)
		})
	}
}
