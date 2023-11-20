package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.redirect.New"

		logger := logger.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		res, err := getter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("alias is not found")
			render.JSON(w, r, response.Error("alias is not found"))
			return
		}
		if err != nil {
			logger.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("some internal error"))
			return
		}

		logger.Info("found url", slog.String("url", res))
		http.Redirect(w, r, res, http.StatusFound)

	}
}
