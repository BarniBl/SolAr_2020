package upload

import (
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
)

type Service interface {
	File(request models.WriteFile) (response models.File, err error)
	Photo(request models.WritePhoto) (response models.Photo, err error)
}

type service struct {
	uploadStorage uploadStorage
	errorWorker   errorWorker
}

func NewService(uploadStorage uploadStorage, errorWorker errorWorker) Service {
	return &service{
		uploadStorage: uploadStorage,
		errorWorker:   errorWorker,
	}
}

func (s *service) File(request models.WriteFile) (response models.File, err error) {
	fileView, err := s.uploadStorage.SaveFile(request)
	if err != nil {
		return
	}

	fileView.ID, err = s.uploadStorage.InsertFile(fileView, request.UserID)
	if err != nil {
		return
	}

	response = fileView

	return
}

func (s *service) Photo(request models.WritePhoto) (response models.Photo, err error) {
	photoView, err := s.uploadStorage.SavePhoto(request)
	if err != nil {
		return
	}

	photoView.ID, err = s.uploadStorage.InsertPhoto(photoView, request.UserID)
	if err != nil {
		return
	}

	response = photoView

	return
}
