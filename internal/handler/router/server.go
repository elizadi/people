package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"people/internal/types"
	"people/internal/usecase"
)

type Server struct {
	usecase *usecase.UseCase
	log     *logrus.Logger
}

func New(uc *usecase.UseCase, log *logrus.Logger) *Server {
	return &Server{
		usecase: uc,
		log:     log,
	}
}

// GetUserInfoBySecondName handler of GET request for retrieving UserInfo by user`s second name
// @Summary Get user details
// @Description Get user information with emails by second name
// @Tags people
//
// @Produce json
// @Param id path int true "User ID"
//
// @Success 200 {object} []types.UserInfo
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id [get]
func (s *Server) GetUserInfoBySecondName(c *gin.Context) {
	name := c.Param("id")

	ctx := context.Background()

	users, err := s.usecase.GetUserInfoBySecondName(ctx, name)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("user not found")
			c.JSON(http.StatusNotFound, types.ErrorResponse{
				Error:   "Not found Error",
				Message: err.Error(),
			})
			return
		}
		s.log.WithError(err).Errorln("Error getting user Info")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": users})
	return
}

// GetAllUsersInfo handler of GET request for retrieving info about all users
// @Summary Get all users details
// @Description Get users information with emails
// @Tags people
//
// @Produce json
//
// @Success 200 {object} []types.UserInfo
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users [get]
func (s *Server) GetAllUsersInfo(c *gin.Context) {
	ctx := context.Background()
	users, err := s.usecase.GetAllUsersInfo(ctx)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Users not found")
			c.JSON(http.StatusNotFound, types.ErrorResponse{
				Error:   "Not found Error",
				Message: err.Error(),
			})
			return
		}
		s.log.WithError(err).Errorln("Error getting users Info")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
	return
}

// GetUserEmails handler of GET request for retrieving all user`s emails
// @Summary Get all user`s emails
// @Description Get all user`s emails
// @Tags people
//
// @Produce json
// @Param id path int true "User ID"
//
// @Success 200 {object} []types.Email
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id/emails [get]
func (s *Server) GetUserEmails(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {

		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	emails, err := s.usecase.GetUserEmails(ctx, idUint)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("User`s emails not found")
			c.JSON(http.StatusNotFound, types.ErrorResponse{
				Error:   "Not found Error",
				Message: err.Error(),
			})
			return
		}
		s.log.WithError(err).Errorln("Error getting user emails")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
	return
}

// GetUserFriends handler of GET request for retrieving all user`s friends
// @Summary Get all user`s friends
// @Description Get all user`s friends
// @Tags people
// @Accept json
//
// @Produce json
// @Param id path int true "User ID"
//
// @Success 200 {object} []types.Friend
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id/friends [get]
func (s *Server) GetUserFriends(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	friends, err := s.usecase.GetUserFriends(ctx, idUint)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("User`s friends found")
			c.JSON(http.StatusNotFound, types.ErrorResponse{
				Error:   "Not found Error",
				Message: err.Error(),
			})
			return
		}
		s.log.WithError(err).Errorln("Error getting user friends")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
	return
}

// CreateUser handler of POST request for add user
// @Summary process POST req for add user
// @Description process POST req for add user
// @Tags people
//
// @Accept json
// @Produce json
// @Param req body types.Name true "first name and second name"
//
// @Success 200 {object} uint64
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users [post]
func (s *Server) CreateUser(c *gin.Context) {
	var name types.Name
	err := c.Bind(&name)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid name")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(name)
	if err != nil {
		s.log.Error("Invalid first name", err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	id, err := s.usecase.CreateUser(ctx, name)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			s.log.WithError(err).Errorln("Nationality not found")
			c.JSON(http.StatusNotFound, types.ErrorResponse{
				Error:   "Not found Error",
				Message: err.Error(),
			})
			return
		}
		s.log.WithError(err).Errorln("Error adding user")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": id})
	return
}

// AddUserEmails handler of POST request for add user`s emails
// @Summary process POST req for add user`s emails
// @Description process POST req for add user`s emails
// @Tags people
//
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param req body types.EmailRequest true "list of user`s emails"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id/emails [post]
func (s *Server) AddUserEmails(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	var emails types.EmailRequest
	err = c.Bind(&emails)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid emails")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(emails)
	if err != nil {
		s.log.Error("Nil emails list", err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()

	err = s.usecase.AddUserEmails(ctx, emails, idUint)
	if err != nil {
		s.log.WithError(err).Errorln("Error adding user`s emails")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "Emails added successfully",
	})
	return
}

// AddUserFriends handler of POST request for add user`s friends
// @Summary process POST req for add user`s friends
// @Description process POST req for add user`s friends
// @Tags people
//
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param req body types.Friends true "list of user`s friends"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id/friends [post]
func (s *Server) AddUserFriends(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	var friends types.Friends
	err = c.Bind(&friends)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid friends info")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}
	validate := validator.New()
	err = validate.Struct(friends)
	if err != nil {
		s.log.Error("Nil list of friends ids", err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	err = s.usecase.AddUserFriends(ctx, friends, idUint)
	if err != nil {
		s.log.WithError(err).Errorln("Error adding user`s friends")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "Friends added successfully",
	})
	return
}

// UpdateUser handler of PUT request for edite user`s info
// @Summary process PUT request for edite user`s info
// @Description process PUT request for edite user`s info
// @Tags people
//
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param req body types.User true "user`s info"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id [put]
func (s *Server) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	var user types.User
	err = c.Bind(&user)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid user info")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		s.log.Error("Invalid first name", err)
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	err = s.usecase.UpdateUser(ctx, user, idUint)
	if err != nil {
		s.log.WithError(err).Errorln("Error updating user")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "User update successfully",
	})
	return
}

// DeleteUser handler of DELETE req for deleting user`s info
// @Summary process DELETE request for deleting user`s info
// @Description process DELETE request for deleting user`s info
// @Tags people
//
// @Accept json
// @Produce json
// @Param id path int true "User ID"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id [delete]
func (s *Server) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		s.log.WithError(err).Errorln("Error getting user id")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	err = s.usecase.DeleteUser(ctx, idUint)
	if err != nil {
		s.log.WithError(err).Errorln("Error deleting user")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "Deleted user successfully",
	})
	return
}

// DeleteEmails handler of DELETE req to delete emails (one or more)
// @Summary process DELETE request to delete emails (one or more)
// @Description process DELETE request to delete emails (one or more)
// @Tags people
//
// @Accept json
// @Produce json
// @Param req body types.EmailIDs true "list email`s ids"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/emails [delete]
func (s *Server) DeleteEmails(c *gin.Context) {
	var emailIDs types.EmailIDs
	err := c.Bind(&emailIDs)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid emails ids")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()

	err = s.usecase.DeleteEmails(ctx, emailIDs.IDs)
	if err != nil {
		s.log.WithError(err).Errorln("Error deleting emails")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "Emails deleted successfully",
	})
	return
}

// DeleteUserFriends handler of DELETE req to delete user`s friendships (one or more)
// @Summary process DELETE request to delete friendships (one or more)
// @Description process DELETE request to delete friendships (one or more)
// @Tags people
//
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param req body types.Friendships true "list of friendships"
//
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/users/:id/friends [delete]
func (s *Server) DeleteUserFriends(c *gin.Context) {
	var friendPairs types.Friendships
	err := c.Bind(&friendPairs)
	if err != nil {
		s.log.WithError(err).Errorln("Invalid friendship info")
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	err = s.usecase.DeleteUserFriends(ctx, friendPairs)
	if err != nil {
		s.log.WithError(err).Errorln("Error deleting user`s friends")
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Error:   "Server Error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponse{
		Message: "Friendships deleted successfully",
	})
	return
}
