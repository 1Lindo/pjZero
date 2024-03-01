package upload

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	resp "main/internal/lib/api/response"
	"main/internal/lib/logger/el"
	"net/http"
)

type Request struct {
	FileName string `json:"file_name"`
	Bytes    []byte `json:"bytes"`
}

type Response struct {
	resp.Response
}

// TODO: fix this
type IConverter interface {
	ConvertToBytes(inputPath string) ([]byte, error)
	ConvertToFile(inputBytes []byte) (string, error)
}

// TODO: fix this
type IPgRepo interface {
	InsertData(fileName string, filePath string) error
	GetData(filePath string) (string, error)
}

func UploadVid(log *slog.Logger, converter IConverter, repo IPgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.video.UploadVid"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("failed to decode request body", el.Err(err))

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//TODO: доработать метод, который получает видео в некой форме + сохраняет го по filePath
		path, err := converter.ConvertToFile(req.Bytes)
		if err != nil {
			log.Error("Failed to convert file", el.Err(err))

			render.JSON(w, r, resp.Error("failed to convert file"))

			return
		}

		err = repo.InsertData(req.FileName, path)
		if err != nil {
			log.Error("Failed to add filepath", el.Err(err))

			render.JSON(w, r, resp.Error("failed to add path"))

			return
		}

		log.Info("filepath added")

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
