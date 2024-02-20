package service

import (
	"log"

	"github.com/sebajax/go-vertical-slice-architecture/internal/user"
	"github.com/sebajax/go-vertical-slice-architecture/pkg/apperror"
)

// CreateUserService interface for DI
type CreateUserService interface {
	CreateUser(user *user.User) (int64, error)
}

// User use cases (port injection)
type createUserService struct {
	userRepository user.UserRepository
}

// Create a new user service use case instance
func NewCreateUserService(repository user.UserRepository) CreateUserService {
	// return the pointer to user service
	return &createUserService{
		userRepository: repository,
	}
}

// Create a new user and store the user in the database
func (service *createUserService) CreateUser(u *user.User) (int64, error) {
	_, check, err := service.userRepository.GetByEmail(u.Email)
	// check if user does not exist and no database error ocurred
	if err != nil {
		// database error
		log.Fatalln(err)
		err := apperror.InternalServerError()
		return 0, err
	}
	if check {
		// user found
		log.Println(u, user.ErrorEmailExists)
		err := apperror.BadRequest(user.ErrorEmailExists)
		return 0, err
	}

	// create the new user and return the id
	userId, err := service.userRepository.Save(u)
	if err != nil {
		// database error
		log.Fatalln(err)
		err := apperror.InternalServerError()
		return 0, err
	}

	// user created successfuly
	return userId, nil
}