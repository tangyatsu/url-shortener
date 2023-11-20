package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

const aliasLength = 9

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"` // тут есть тег для пакета validator/v10 который проверяет валидность url
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func New(logger *slog.Logger, saver URLSaver) http.HandlerFunc { // к нам будет приходить POST запрос с джейсоном с ссылкой для сохранения
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		logger = logger.With(slog.String("fn", fn), slog.String("request_id", middleware.GetReqID(r.Context())))

		var body Request

		err := render.DecodeJSON(r.Body, &body)
		if err != nil {
			logger.Error("failed to decode req body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode req body")) // часть chi для удобной работы

			return
		}

		logger.Info("body decoded", slog.Any("request", body))

		if err = validator.New().Struct(body); err != nil {

			logger.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		alias := body.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength) //TODO надо прочекать что алиас уникальный (мб несколько раз перегененировать)
		}

		id, err := saver.SaveURL(body.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Info("url already exists", slog.String("url", body.URL)) // чет навертели эта ошибка будет тупо если алиас совпадает, но такого урла мб еще нет

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			logger.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		logger.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{Response: response.OK(), Alias: alias})
	}
}
