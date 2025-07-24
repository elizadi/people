package usecase

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"people/internal/repository/enrichment"
	"people/internal/repository/storage"
	"people/internal/types"
)

type UseCase struct {
	storage    *storage.Storage
	enrichment *enrichment.Enrichment
	log        *logrus.Logger
}

func New(storage *storage.Storage, enrichment *enrichment.Enrichment, log *logrus.Logger) *UseCase {
	return &UseCase{
		storage:    storage,
		enrichment: enrichment,
		log:        log,
	}
}

func (s *UseCase) GetUserInfoByID(ctx context.Context, id uint64) (types.UserInfo, error) {
	user, err := s.storage.GetUserInfoByID(ctx, id)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Not found user info")
			return types.UserInfo{}, types.ErrNotFound
		}
		s.log.WithError(err).Errorln("Can`t get user info")
		return types.UserInfo{}, err
	}

	return user, nil
}

func (s *UseCase) GetAllUsersInfo(ctx context.Context) ([]types.UserInfo, error) {
	users, err := s.storage.GetAllUsersInfo(ctx)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Not found all users info")
			return []types.UserInfo{}, types.ErrNotFound
		}
		s.log.WithError(err).Errorln("Can`t get all users info")
		return []types.UserInfo{}, err
	}

	return users, nil
}

func (s *UseCase) GetAllUserEmails(ctx context.Context, id uint64) ([]types.Email, error) {
	emails, err := s.storage.GetAllUserEmails(ctx, id)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Not found user`s emails")
			return []types.Email{}, types.ErrNotFound
		}
		s.log.WithError(err).Errorln("Can`t get all user`s emails")
		return []types.Email{}, err
	}

	return emails, nil
}

func (s *UseCase) GetUserFriends(ctx context.Context, id uint64) ([]types.Friend, error) {
	friends, err := s.storage.GetUserFriends(ctx, id)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Not found user`s friends")
			return []types.Friend{}, types.ErrNotFound
		}
		s.log.WithError(err).Errorln("Can`t get user`s friends")
		return []types.Friend{}, err
	}

	return friends, nil
}

func (s *UseCase) CreateUser(ctx context.Context, fullName types.Name) (uint64, error) {
	age, err := s.enrichment.Age(fullName.FirstName)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t get user age")
		return 0, err
	}

	gender, err := s.enrichment.Gender(fullName.FirstName)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t get user gender")
		return 0, err
	}

	nationality, err := s.enrichment.Nationality(fullName.FirstName)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t get user nationality")
		return 0, err
	}

	user := types.User{
		Name: types.Name{
			FirstName: fullName.FirstName,
			LastName:  fullName.LastName,
		},
		Gender:      gender,
		Nationality: nationality,
		Age:         age,
	}

	id, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t add user")
		return 0, err
	}

	return id, nil
}

// AddUserEmails - can add one or more user`s emails
func (s *UseCase) AddUserEmails(ctx context.Context, emails types.EmailRequest, id uint64) error {
	err := s.storage.AddUserEmails(ctx, emails, id)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t add user`s emails")
	}
	return err
}

// AddUserFriends - can add one or more user`s friends
func (s *UseCase) AddUserFriends(ctx context.Context, friends types.Friends, userID uint64) error {
	err := s.storage.AddUserFriends(ctx, friends, userID)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t add user friends")
	}
	return err
}

func (s *UseCase) UpdateUser(ctx context.Context, user types.User, id uint64) error {
	err := s.storage.UpdateUser(ctx, user, id)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t update user")
	}
	return err
}

func (s *UseCase) DeleteUser(ctx context.Context, id uint64) error {
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t delete user")
	}
	return err
}

// DeleteEmails - can delete one or more emails
func (s *UseCase) DeleteEmails(ctx context.Context, emails []uint64) error {
	err := s.storage.DeleteEmails(ctx, emails)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t delete emails")
	}
	return err
}

// DeleteUserFriends - can delete one or more user`s friends
func (s *UseCase) DeleteUserFriends(ctx context.Context, friendsPairs types.Friendships) error {
	err := s.storage.DeleteUserFriends(ctx, friendsPairs)
	if err != nil {
		s.log.WithError(err).Errorln("Can`t delete user friends")
	}
	return err
}
