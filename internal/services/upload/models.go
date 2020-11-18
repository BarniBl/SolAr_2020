package upload

import (
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
)

type uploadStorage interface {
	SaveFile(file models.WriteFile) (fileView models.File, err error)
	SavePhoto(file models.WritePhoto) (fileView models.Photo, err error)
	//SaveFiles(files []models.WriteFile) (fileView []models.File, err error)

	InsertFile(file models.File, userID int) (fileID int, err error)
	InsertPhoto(photo models.Photo, userID int) (photoID int, err error)
}

type errorWorker interface {
	NewError(httpCode int, responseError error, fullError error) (err error)
}
