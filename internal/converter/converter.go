package converter

import (
	"fmt"
	"os"
)

type IConverter interface {
	ConvertToBytes(filePath string) ([]byte, error)
	ConvertToFile(fileName string, path string, inputBytes []byte) (string, error)
}

type Converter struct {
}

// TODO: рализация метода, который находит файл по filePath + конвертирует его в необходимый формат
func (c *Converter) ConvertToBytes(filePath string) ([]byte, error) {
	const op string = "internal.converter.converter.ConvertToBytes" //op == operation
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return file, nil
}

// TODO: рализация метода, который получает []byte + cоздает файл для записи + записывает в этот файл
func (c *Converter) ConvertToFile(fileName string, path string, inputBytes []byte) (string, error) {
	const op string = "internal.converter.converter.ConvertToFile"

	//cоздаю временный файл

	fileOnDisk, err := os.Create(fmt.Sprintf("%s/%s", path, fileName))
	if err != nil {
		return "", err
	}
	defer fileOnDisk.Close()

	_, err = fileOnDisk.Write(inputBytes)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Sprintf(op + "success"), nil

}
