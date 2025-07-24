package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/sirupsen/logrus"
	"people/internal/types"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *logger.Entry
}

func (s *Storage) Migrations(ctx context.Context) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Could not ping the db!")
		return err
	}

	_, err = connection.Exec(ctx, createUsersTableTemplate)
	if err != nil {
		s.logger.WithError(err).Errorln("Error create Users table")
		return err
	}

	_, err = connection.Exec(ctx, createEmailsTableTemplate)
	if err != nil {
		s.logger.WithError(err).Errorln("Error create Emails table")
		return err
	}

	_, err = connection.Exec(ctx, createFriendsTableTemplate)
	if err != nil {
		s.logger.WithError(err).Errorln("Error create Friends table")
		return err
	}
	_, err = connection.Exec(ctx, createFriendsIndexTemplates)
	if err != nil {
		s.logger.WithError(err).Errorln("Error create Friends additional index")
		return err
	}

	return nil
}

func New(ctx context.Context, dbUrl string, log *logger.Logger) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return &Storage{}, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.WithError(err).Errorln("Could not ping the db!")
		return &Storage{}, err
	}

	logEntry := log.WithField("package", "storage")

	r := &Storage{pool: pool, logger: logEntry}
	err = r.Migrations(ctx)
	if err != nil {
		return &Storage{}, err
	}
	return r, nil
}

func (s *Storage) GetAllUsersInfo(ctx context.Context) ([]types.UserInfo, error) {
	connection, err := s.pool.Acquire(ctx)
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return []types.UserInfo{}, err
	}

	defer connection.Release()

	rows, err := connection.Query(ctx, GetAllUsersTemplate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.WithError(err).Errorln("No such rows in Users")
			return nil, types.ErrNotFound
		}
		s.logger.WithError(err).Errorln("Error getting all users info")
		return []types.UserInfo{}, err
	}

	var users []types.UserInfo
	var errs []error

	for rows.Next() {
		var user types.UserInfo
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Gender,
			&user.Age,
			&user.Nationality,
			&user.Emails,
		)

		if err != nil {
			s.logger.WithError(err).Errorln("Error getting user info")
			errs = append(errs, err)
		}

		users = append(users, user)
	}

	err = errors.Join(errs...)
	if err != nil {
		s.logger.WithError(err).Errorln("Error getting all users info")
		return []types.UserInfo{}, err
	}

	return users, nil
}

func (s *Storage) GetUserInfoByID(ctx context.Context, id uint64) (types.UserInfo, error) {
	connection, err := s.pool.Acquire(ctx)
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return types.UserInfo{}, err
	}

	defer connection.Release()

	var user types.UserInfo

	err = connection.QueryRow(ctx, GetUserAllInfoTemplate, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Gender,
		&user.Age,
		&user.Nationality,
		&user.Emails,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.WithError(err).Errorln("No such row in Users")
			return types.UserInfo{}, types.ErrNotFound
		}
		s.logger.WithError(err).Errorln("Error getting user info")
		return types.UserInfo{}, err
	}

	return user, nil
}

func (s *Storage) GetAllUserEmails(ctx context.Context, id uint64) ([]types.Email, error) {
	connection, err := s.pool.Acquire(ctx)
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return []types.Email{}, err
	}

	defer connection.Release()

	rows, err := connection.Query(ctx, GetAllUserEmailsTemplate, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.WithError(err).Errorln("No such rows in Email")
			return []types.Email{}, types.ErrNotFound
		}
		s.logger.WithError(err).Errorln("Error getting user`s emails")
		return []types.Email{}, err
	}

	var emails []types.Email
	var errs []error

	for rows.Next() {
		var email types.Email
		err = rows.Scan(
			&email.ID,
			&email.Email,
		)

		if err != nil {
			s.logger.WithError(err).Errorln("Error getting email")
			errs = append(errs, err)
		}

		emails = append(emails, email)
	}

	err = errors.Join(errs...)
	if err != nil {
		s.logger.WithError(err).Errorln("Error getting emails")
		return []types.Email{}, err
	}

	return emails, nil
}

func (s *Storage) GetUserFriends(ctx context.Context, id uint64) ([]types.Friend, error) {
	connection, err := s.pool.Acquire(ctx)
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return []types.Friend{}, err
	}

	defer connection.Release()

	rows, err := connection.Query(ctx, GetUserFriendsTemplate, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.WithError(err).Errorln("No such rows in Friends")
			return []types.Friend{}, types.ErrNotFound
		}
		s.logger.WithError(err).Errorln("Error getting all friends info")
		return []types.Friend{}, err
	}

	var friends []types.Friend
	var errs []error

	for rows.Next() {
		var friend types.Friend
		err = rows.Scan(
			&friend.FriendID,
			&friend.FirstName,
			&friend.LastName,
		)

		if err != nil {
			s.logger.WithError(err).Errorln("Error getting friend name")
			errs = append(errs, err)
		}

		friends = append(friends, friend)
	}

	err = errors.Join(errs...)
	if err != nil {
		s.logger.WithError(err).Errorln("Error getting friends")
		return []types.Friend{}, err
	}

	return friends, nil
}

func (s *Storage) CreateUser(ctx context.Context, user types.User) (uint64, error) {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return 0, err
	}

	defer connection.Release()

	var id uint64

	err = connection.QueryRow(
		ctx,
		AddUserInfoTemplate,
		user.FirstName,
		user.LastName,
		user.Gender,
		user.Nationality,
		user.Age,
	).Scan(&id)
	if err != nil {
		s.logger.WithError(err).Errorln("Failed to add user")
		return 0, err
	}
	return id, nil
}

// AddUserEmails - can add one or more user`s emails
func (s *Storage) AddUserEmails(ctx context.Context, emails types.EmailRequest, id uint64) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	if len(emails.Emails) == 0 {
		s.logger.Errorln("No emails found!")
		return types.ErrNotFound
	}

	batch := &pgx.Batch{}

	for _, email := range emails.Emails {
		batch.Queue(
			AddEmailTemplate,
			id,
			strings.TrimSpace(email),
		)
	}

	err = connection.SendBatch(ctx, batch).Close()
	if err != nil {
		s.logger.WithError(err).Errorln("Error sending batch!")
		return err
	}

	return nil
}

// AddUserFriends - can add one or more user`s friends
func (s *Storage) AddUserFriends(ctx context.Context, friends types.Friends, userID uint64) error {
	connection, err := s.pool.Acquire(ctx)
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	if len(friends.FriendsIDs) == 0 {
		s.logger.Errorln("No fiends found!")
		return errors.New("No friends")
	}

	batch := &pgx.Batch{}

	for _, friend := range friends.FriendsIDs {
		if friend == userID {
			s.logger.Warnf("Attempted to add self as friend (userID: %d)\n", userID)
			continue
		}
		// add ids in ascending order to avoid duplicating pairs (since the pairs of id [1, 2] and [2, 1] equal
		if userID < friend {
			batch.Queue(
				AddFriendshipTemplate,
				userID,
				friend,
			)
		} else {
			batch.Queue(
				AddFriendshipTemplate,
				friend,
				userID,
			)
		}
	}

	err = connection.SendBatch(ctx, batch).Close()
	if err != nil {
		s.logger.WithError(err).Errorln("Error sending batch!")
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user types.User, id uint64) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	commandTag, err := connection.Exec(
		ctx,
		UpdateUserInfoTemplate,
		id,
		user.FirstName,
		user.LastName,
		user.Gender,
		user.Nationality,
		user.Age,
	)
	if err != nil {
		s.logger.WithError(err).Errorln("Failed to update user")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		err = fmt.Errorf("Not found user with id %d", id)
		s.logger.Errorln(err)
		return err
	}
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id uint64) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	commandTag, err := connection.Exec(
		ctx,
		DeleteUserTemplate,
		id,
	)
	if err != nil {
		s.logger.WithError(err).Errorln("Failed to delete user")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		err = fmt.Errorf("Not found user with id %d", id)
		s.logger.Errorln(err)
		return err
	}
	return nil
}

// DeleteEmails - can delete one or more emails
func (s *Storage) DeleteEmails(ctx context.Context, emails []uint64) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	if len(emails) == 0 {
		s.logger.Errorln("No emails found!")
		return errors.New("No emails")
	}

	batch := &pgx.Batch{}

	for _, email := range emails {
		batch.Queue(
			DeleteEmailTemplate,
			email,
		)
	}

	err = connection.SendBatch(ctx, batch).Close()
	if err != nil {
		s.logger.WithError(err).Errorln("Error sending batch!")
		return err
	}

	return nil
}

// DeleteUserFriends - can delete one or more user`s friends
func (s *Storage) DeleteUserFriends(ctx context.Context, friendsPairs types.Friendships) error {
	connection, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.WithError(err).Errorln("Error whole acquiring connection from the database pool!")
		return err
	}

	defer connection.Release()

	if len(friendsPairs.Friends) == 0 {
		s.logger.Errorln("No fiends found!")
		return errors.New("No friends")
	}

	batch := &pgx.Batch{}

	for _, friends := range friendsPairs.Friends {
		// use ids in ascending order to avoid duplicating pairs (since the pairs of id [1, 2] and [2, 1] equal
		if friends.IDFirstUser < friends.IDSecondUser {
			batch.Queue(
				DeleteFriendshipTemplate,
				friends.IDFirstUser,
				friends.IDSecondUser,
			)
		} else {
			batch.Queue(
				DeleteFriendshipTemplate,
				friends.IDSecondUser,
				friends.IDFirstUser,
			)
		}
	}

	err = connection.SendBatch(ctx, batch).Close()
	if err != nil {
		s.logger.WithError(err).Errorln("Error sending batch!")
		return err
	}

	return nil
}
