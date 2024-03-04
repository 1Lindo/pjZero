package download

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	resp "main/internal/lib/api/response"
	"main/internal/lib/logger/el"
	"net/http"
	"strconv"
)

type Request struct {
	FileName string `json:"file_name"`
}

type Response struct {
	resp.Response
	FilePath  string `json:"filePath"`
	FileBytes []byte `json:"fileBytes"`
}

// TODO: fix this
type IConverter interface {
	ConvertToBytes(inputPath string) ([]byte, error)
	ConvertToFile(fileName string, path string, inputBytes []byte) (string, error)
}

// TODO: fix this
type IPgRepo interface {
	InsertData(fileName string, filePath string) error
	GetData(fileName string) (string, error)
}

func DownloadVid(log *slog.Logger, converter IConverter, repo IPgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.http-server.handlers.video.DownloadVid"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		fileName := r.URL.Query().Get("file_name")

		log.Info("request body decoded", slog.Any("request", fileName))

		path, err := repo.GetData(fileName)
		if err != nil {
			log.Error("Failed to get filepath", el.Err(err))

			render.JSON(w, r, resp.Error("failed to get path"))

			return
		}

		log.Info("filepath received:", path)

		//TODO: метод, который получает filePath + находит нужный файл + передает его клиенту как []bytes
		//TODO: иправить путь к файлу, тк не получается найти файл используя переменную path
		fileBytes, err := converter.ConvertToBytes("/Users/leoniddomanin/golangprojects/pjZero/cmd/pjZero/vidoes/12/12.mp4")
		if err != nil {
			//TODO: обработать ошибку
		}

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
		w.Header().Set("Path", path)
		w.WriteHeader(http.StatusOK)
		w.Write(fileBytes)
	}
}
