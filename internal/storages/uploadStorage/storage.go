package uploadStorage

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	SaveFile(file models.WriteFile) (fileView models.File, err error)
	SavePhoto(file models.WritePhoto) (fileView models.Photo, err error)

	InsertFile(file models.File, userID int) (fileID int, err error)
	InsertPhoto(photo models.Photo, userID int) (photoID int, err error)

	SelectCountFiles(fileIDs []int, userID int) (countFiles int, err error)
	SelectCountPhotos(photoIDs []int, userID int) (countPhotos int, err error)

	SelectFiles(fileIDs []int) (files map[int]models.File, err error)
	SelectPhotos(photoIDs []int) (photos map[int]models.Photo, err error)
}

type storage struct {
	photoPath string
	filePath  string
	db        *sql.DB
}

func NewStorage(photoPath, filePath string, db *sql.DB) Storage {
	return &storage{
		photoPath: photoPath,
		filePath:  filePath,
		db:        db,
	}
}

func (s *storage) SaveFile(file models.WriteFile) (fileView models.File, err error) {
	fileName, err := s.createFileName(file.Name)
	if err != nil {
		return
	}

	filePath, err := s.createFilePath(file.Name)
	if err != nil {
		return
	}

	fileView.Name = file.Name
	fileView.URL = s.filePath + "/" + filePath + "/" + fileName

	if err = os.MkdirAll(s.filePath+"/"+filePath, 0777); err != nil {
		fmt.Println("Error on MkdirAll: ", err)
		return
	}

	writeFile, err := os.Create(fileView.URL)
	if err != nil {
		fmt.Println("Cannot create file: ", err)
		return
	}
	defer writeFile.Close()

	readFile, err := file.File.Open()
	if err != nil {
		fmt.Println("Cannot open file: ", err)
		return
	}
	defer readFile.Close()

	n, err := io.Copy(writeFile, readFile)
	fmt.Println("File copied: err=", err, ", n= ", n)

	return
}

func (s *storage) SavePhoto(photo models.WritePhoto) (photoView models.Photo, err error) {
	fileName, err := s.createFileName(photo.Name)
	if err != nil {
		return
	}

	filePath, err := s.createFilePath(fileName)
	if err != nil {
		return
	}
	photoView.Name = photo.Name
	photoView.URL = s.photoPath + "/" + filePath + "/" + fileName

	if err = os.MkdirAll(s.photoPath+"/"+filePath, 0777); err != nil {
		return
	}

	writeFile, err := os.Create(photoView.URL)
	if err != nil {
		return
	}
	defer writeFile.Close()

	readFile, err := photo.File.Open()
	if err != nil {
		return
	}
	defer readFile.Close()

	_, err = io.Copy(writeFile, readFile)

	return
}

func (s *storage) InsertFile(file models.File, userID int) (fileID int, err error) {
	const sqlQuery = `
	INSERT INTO files (title, url, user_id)
	VALUES ($1, $2, $3)
	RETURNING id;`

	err = s.db.QueryRow(sqlQuery, file.Name, file.URL, userID).Scan(&fileID)
	return
}

func (s *storage) InsertPhoto(photo models.Photo, userID int) (photoID int, err error) {
	const sqlQuery = `
	INSERT INTO photos (title, url, user_id)
	VALUES ($1, $2, $3)
	RETURNING id;`

	err = s.db.QueryRow(sqlQuery, photo.Name, photo.URL, userID).Scan(&photoID)
	return
}

func (s *storage) createFilePath(name string) (pathName string, err error) {
	if len(name) < 3 {
		return pathName, errors.New("Некорректное имя пути")
	}

	return name[:2], nil
}

func (s *storage) createFileName(name string) (fileName string, err error) {
	postfix, err := s.extractFormatFile(name)
	if err != nil {
		return
	}

	h := md5.New()
	h.Write([]byte(time.Now().String() + name))

	fileName = fmt.Sprintf("%x", h.Sum(nil))

	return fileName + "." + postfix, nil
}

func (s *storage) extractFormatFile(fileName string) (postfix string, err error) {
	if fileName == "" {
		err = fmt.Errorf("filename must not be empty")
		return
	}
	parts := strings.Split(fileName, ".")
	if len(parts) == 0 {
		return postfix, errors.New("no wanna deal with non-postfix files")
	}
	postfix = parts[len(parts)-1]
	return
}

func (s *storage) SelectCountFiles(fileIDs []int, userID int) (countFiles int, err error) {
	if len(fileIDs) == 0 {
		return 0, err
	}

	const sqlQueryTemplate = `
	SELECT count(*)
	FROM upload.files AS f
	WHERE f.user_id = ? AND f.id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(fileIDs))

	var params []interface{}
	params = append(params, userID)
	for i, _ := range fileIDs {
		params = append(params, fileIDs[i])
	}

	for i := 1; i <= len(fileIDs)*1+1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	err = s.db.QueryRow(sqlQuery, params...).Scan(&countFiles)

	return
}

func (s *storage) SelectCountPhotos(photoIDs []int, userID int) (countPhotos int, err error) {
	if len(photoIDs) == 0 {
		return 0, err
	}

	const sqlQueryTemplate = `
	SELECT count(*)
	FROM photos AS p
	WHERE p.user_id = ? AND p.id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(photoIDs))

	var params []interface{}
	params = append(params, userID)
	for i, _ := range photoIDs {
		params = append(params, photoIDs[i])
	}

	for i := 1; i <= len(photoIDs)*1+1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	err = s.db.QueryRow(sqlQuery, params...).Scan(&countPhotos)

	return
}

func createIN(count int) (queryIN string) {
	queryIN = "("
	for i := 0; i < count; i++ {
		queryIN += "?, "
	}
	queryINRune := []rune(queryIN)
	queryIN = string(queryINRune[:len(queryINRune)-2])
	queryIN += ")"
	return
}

func (s *storage) SelectFiles(fileIDs []int) (files map[int]models.File, err error) {
	files = make(map[int]models.File)
	if len(fileIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT f.id, f.title, f.url
	FROM files AS f
	WHERE f.id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(fileIDs))

	var params []interface{}

	for i, _ := range fileIDs {
		params = append(params, fileIDs[i])
	}

	for i := 1; i <= len(fileIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempFile models.File
		err = rows.Scan(&tempFile.ID, &tempFile.Name, &tempFile.URL)
		if err != nil {
			return
		}
		files[tempFile.ID] = tempFile
	}

	return
}

func (s *storage) SelectPhotos(photoIDs []int) (photos map[int]models.Photo, err error) {
	photos = make(map[int]models.Photo)
	if len(photoIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT p.id, p.title, p.url
	FROM photos AS p
	WHERE p.id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(photoIDs))

	var params []interface{}

	for i, _ := range photoIDs {
		params = append(params, photoIDs[i])
	}

	for i := 1; i <= len(photoIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempPhoto models.Photo
		err = rows.Scan(&tempPhoto.ID, &tempPhoto.Name, &tempPhoto.URL)
		if err != nil {
			return
		}
		photos[tempPhoto.ID] = tempPhoto
	}

	return
}

//func (s *storage) SaveFiles(files []models.WriteFile) (fileViews []models.File, err error) {
//	for i, _ := range files {
//		fileView, err := s.SaveFile(files[i])
//		if err != nil {
//			return
//		}
//		fileViews = append(fileViews, fileView)
//	}
//
//	return
//}
