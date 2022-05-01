package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterInputUser) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, filelocation string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterInputUser) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"
	NewUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}

	return NewUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	user, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	println(user.PasswordHash)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	user, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(id int, filelocation string) (User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = filelocation
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}
