package types

import "errors"

var ErrNotFound = errors.New("Not found")

type User struct {
	Name
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Age         uint8  `json:"age"`
}

type Name struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"omitempty"`
}

type UserInfo struct {
	ID uint64 `json:"id"`
	User
	Emails []string `json:"emails"`
}

type Email struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
}

type EmailRequest struct {
	Emails []string `json:"emails" validate:"required"`
}

type EmailIDs struct {
	IDs []uint64 `json:"ids"`
}

type Friends struct {
	FriendsIDs []uint64 `json:"friends_ids" validate:"required"`
}

type Friendship struct {
	IDFirstUser  uint64 `json:"id_first_friend" validate:"required"`
	IDSecondUser uint64 `json:"id_second_friend" validate:"required"`
}

type Friendships struct {
	Friends []Friendship
}

type Friend struct {
	FriendID uint64 `json:"friend_id"`
	Name
}
