package posts

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	interviewModels "github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"sort"
	"strconv"
)

type Service interface {
	Create(request models.InputPost) (response models.Post, err error)
	GetList(request models.GetPostListRequest) (response []models.PostResult, err error)
	SetMark(request models.MarkPost) (err error)
	Delete(request models.DeletePostRequest) (err error)
}

type service struct {
	postsStorage     postStorage
	uploadStorage    uploadStorage
	interviewStorage interviewStorage
	paymentClient    paymentClient
	groupClient      groupClient
	accountClient    accountClient
	errorWorker      errorWorker
}

func NewService(postsStorage postStorage, uploadStorage uploadStorage, interviewStorage interviewStorage, groupClient groupClient, accountClient accountClient, paymentClient paymentClient, errorWorker errorWorker) Service {
	return &service{
		postsStorage:     postsStorage,
		uploadStorage:    uploadStorage,
		interviewStorage: interviewStorage,
		groupClient:      groupClient,
		accountClient:    accountClient,
		paymentClient:    paymentClient,
		errorWorker:      errorWorker,
	}
}

func (s *service) Create(request models.InputPost) (response models.Post, err error) {
	if err = s.validateCreate(request); err != nil {
		return
	}

	err = s.groupClient.CheckPermission(request.CreateBy, request.GroupID, CreatePostActionID)
	if err != nil {
		return response, err
	}

	if err = s.checkFiles(request.Files, request.CreateBy); err != nil {
		return
	}

	if err = s.checkPhotos(request.Photos, request.CreateBy); err != nil {
		return
	}

	response.ID, err = s.postsStorage.InsertPost(request)
	if err != nil {
		return
	}

	err = s.interviewStorage.InsertInterviews(request.Interviews, response.ID)
	if err != nil {
		return
	}

	createRequest := models.CreateRequest{
		CreateBy: request.CreateBy,
		GroupID:  request.GroupID,
		PostID:   response.ID,
		Payments: request.Payments,
	}

	_, err = s.paymentClient.Create(createRequest)
	if err != nil {
		return
	}

	// TODO CHANGE TO CONST
	err = s.postsStorage.UpdatePostStatus(response.ID, request.GroupID, PostStatusCreated)
	if err != nil {
		return
	}

	return
}

func (s *service) validateCreate(post models.InputPost) (err error) {
	if len(post.Files) > FilesLimit {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorFilesLimit, errors.Wrap(ErrorFilesLimit, strconv.Itoa(len(post.Files))))
	}

	if len(post.Photos) > PhotosLimit {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorPhotosLimit, errors.Wrap(ErrorFilesLimit, strconv.Itoa(len(post.Photos))))
	}

	if len(post.Payments) > PaymentsLimit {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorPaymentsLimit, errors.Wrap(ErrorFilesLimit, strconv.Itoa(len(post.Payments))))
	}

	if len(post.Interviews) > InterviewsLimit {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorInterviewsLimit, errors.Wrap(ErrorFilesLimit, strconv.Itoa(len(post.Interviews))))
	}

	return
}

func (s *service) checkFiles(fileIDs []int, userID int) (err error) {
	countFiles, err := s.uploadStorage.SelectCountFiles(fileIDs, userID)
	if err != nil {
		return s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	if countFiles != len(fileIDs) {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorFilesNotFound, ErrorFilesNotFound)
	}

	return
}

func (s *service) checkPhotos(photoIDs []int, userID int) (err error) {
	countFiles, err := s.uploadStorage.SelectCountPhotos(photoIDs, userID)
	if err != nil {
		return
	}

	if countFiles != len(photoIDs) {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorPhotosNotFound, ErrorPhotosNotFound)
	}

	return
}

func (s *service) GetList(request models.GetPostListRequest) (response []models.PostResult, err error) {
	err = s.groupClient.CheckPermission(request.UserID, request.GroupID, GetPostActionID)
	if err != nil {
		return response, err
	}

	posts, err := s.postsStorage.SelectPosts(request)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	if len(posts) == 0 {
		response = make([]models.PostResult, 0)
		return
	}

	postsMap := make(map[int]models.PostResult)
	for index, post := range posts {
		postsMap[post.ID] = models.PostResult{
			ID:          post.ID,
			CreateBy:    post.CreateBy,
			CreatAt:     post.CreatAt,
			PublishDate: post.PublishDate,
			GroupID:     post.GroupID,
			Text:        post.Text,
			Status:      post.Status,
			Photos:      make([]models.Photo, 0),
			Files:       make([]models.File, 0),
			Interviews:  make([]interviewModels.InterviewResult, 0),
			Payments:    make([]models.Payment, 0),
			Order:       index,
			Marked:      post.Marked,
		}
	}

	postIDs := make([]int, 0)
	for i, _ := range posts {
		postIDs = append(postIDs, posts[i].ID)
	}

	interviews, err := s.interviewStorage.SelectInterviewsResults(postIDs, request.UserID)
	if err != nil {
		return
	}

	payments, err := s.paymentClient.GetByPostIDs(postIDs)
	if err != nil {
		return
	}

	matchPostPhoto, err := s.postsStorage.SelectPhotoIDs(postIDs)
	if err != nil {
		return
	}

	photoIDs := make([]int, 0)
	for i, _ := range matchPostPhoto {
		photoIDs = append(photoIDs, matchPostPhoto[i].PhotoID)
	}

	matchPostFile, err := s.postsStorage.SelectFileIDs(postIDs)
	if err != nil {
		return
	}

	fileIDs := make([]int, 0)
	for i, _ := range matchPostFile {
		fileIDs = append(fileIDs, matchPostFile[i].FileID)
	}

	photos, err := s.uploadStorage.SelectPhotos(photoIDs)
	if err != nil {
		return
	}

	files, err := s.uploadStorage.SelectFiles(fileIDs)
	if err != nil {
		return
	}

	for _, interview := range interviews {
		tempPost := postsMap[interview.PostID]
		tempPost.Interviews = append(tempPost.Interviews, interview)
		postsMap[interview.PostID] = tempPost
	}

	for _, payment := range payments {
		tempPost := postsMap[payment.PostID]
		tempPost.Payments = append(tempPost.Payments, payment)
		postsMap[payment.PostID] = tempPost
	}

	for _, match := range matchPostPhoto {
		tempPost := postsMap[match.PostID]
		tempPost.Photos = append(tempPost.Photos, photos[match.PhotoID])
		postsMap[match.PostID] = tempPost
	}

	for _, match := range matchPostFile {
		tempPost := postsMap[match.PostID]
		tempPost.Files = append(tempPost.Files, files[match.FileID])
		postsMap[match.PostID] = tempPost
	}

	for _, post := range postsMap {
		response = append(response, post)
	}

	sortPost := models.Posts{Posts: response}
	sort.Sort(&sortPost)

	for i, _ := range sortPost.Posts {
		var user account.User
		user, err = s.accountClient.GetUserByUid(sortPost.Posts[i].CreateBy)
		if err != nil {
			return
		}
		sortPost.Posts[i].Author = user
	}

	return sortPost.Posts, nil
}

func (s *service) SetMark(request models.MarkPost) (err error) {
	err = s.groupClient.CheckPermission(request.UserID, request.GroupID, MarkPostActionID)
	if err != nil {
		return err
	}

	err = s.postsStorage.SetMark(request.PostID, request.Mark, request.GroupID)
	if err != nil {
		err = errors.New("cannot set mark: " + err.Error())
	}
	return
}

func (s *service) Delete(request models.DeletePostRequest) (err error) {
	err = s.groupClient.CheckPermission(request.UserID, request.GroupID, DeletePostActionID)
	if err != nil {
		return err
	}

	err = s.postsStorage.UpdatePostStatus(request.PostID, request.GroupID, PostStatusRemoved)
	if err != nil {
		err = errors.New("cannot remove post: " + err.Error())
	}
	return
}
