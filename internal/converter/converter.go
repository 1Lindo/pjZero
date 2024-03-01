package converter

import (
	"bytes"
	"fmt"
	"os/exec"
)

type IConverter interface {
	ConvertToBytes(filePath string) ([]byte, error)
	ConvertToFile(inputBytes []byte) (string, error)
}

type Converter struct {
}

// TODO: рализация метода, который находит файл по filePath + конвертирует его в необходимый формат
func (c *Converter) ConvertToBytes(filePath string) ([]byte, error) {
	const op string = "internal.converter.converter.ConvertToBytes" //op == operation
	cmd := exec.Command("ffmpeg", "-i", filePath, "-f", "rawvideo", "-")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	videoBytes := out.Bytes()

	return videoBytes, nil
}

// TODO: рализация метода, который получает видео в неком формате + конвертирует его в необходимый формат + отдает на выходе filePath
func (c *Converter) ConvertToFile(inputBytes []byte) (string, error) {
	const op string = "internal.converter.converter.ConvertToFile"

	// Получение имени директории для сохранения файлов
	//outputDir := "/Users/USERDIR/Desktop/videos"
	//
	//// Создание директории, если её нет
	//err := os.Mkdir(outputDir, os.ModePerm)
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	//cоздаю временный файл
	//tempFile := "/Users/USERDIR/Desktop/videos/temp_file.mp4"
	//err := os.WriteFile(tempFile, inputBytes, 0644)
	//if err != nil {
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}

	//tempFile, err := os.CreateTemp("", "video.mp4")
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	//defer os.Remove(tempFile.Name())
	//defer tempFile.Close()

	//_, err = tempFile.Write(inputBytes)
	//if err != nil {
	//	return fmt.Errorf("%s: %w", op, err)
	//}

	//outputFilePath := "/Users/USERDIR/Desktop/videos/output.mp4"
	//os.Create(outputFilePath)
	//cmd := exec.Command("ffmpeg", "-i", tempFile, outputFilePath)
	//err = cmd.Run()
	//if err != nil {
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}

	return op, nil

}
