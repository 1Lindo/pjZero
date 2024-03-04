package upload

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	resp "main/internal/lib/api/response"
	"main/internal/lib/logger/el"
	"net/http"
	"os"
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
	ConvertToFile(fileName string, path string, inputBytes []byte) (string, error)
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

		//videoID, _ := r.Context().Value("videoId").(string)
		videoID := "12"

		//if len(videoID) == 0 {
		//	videoID = "default"
		//}
		err := os.MkdirAll(fmt.Sprintf("./vidoes/%s", videoID), os.ModePerm) // создаем директорию с id продукта, если еще не создана
		if err != nil {
			//TODO: обработать ошибку
		}

		path := fmt.Sprintf("./vidoes/%s", videoID)

		// извлечение файла из параметров запроса
		//form := r.MultipartForm
		var fileName string = "file12"
		videExt := "mp4"
		//for key := range form.File {
		//	fileName = key
		//	arr := strings.Split(fileName, ".")
		//	if len(arr) > 1 {
		//		videExt = arr[len(arr)-1]
		//	}
		//	continue
		//}
		//извлекаем из содержимое файла
		file, _, err := r.FormFile(fileName)
		if err != nil {
			//TODO: обработать ошибку
		}
		defer file.Close()

		var req Request
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("failed to decode request body", el.Err(err))

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		fullFileName := fmt.Sprintf("%s.%s", videoID, videExt)

		log.Info("request body decoded", slog.Any("request", req))

		//TODO: доработать метод, который получает видео в некой форме + сохраняет го по filePath
		_, err = converter.ConvertToFile(fullFileName, path, fileBytes)
		if err != nil {
			log.Error("Failed to convert file", el.Err(err))

			render.JSON(w, r, resp.Error("failed to convert file"))

			return
		}

		err = repo.InsertData(fullFileName, path)
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
