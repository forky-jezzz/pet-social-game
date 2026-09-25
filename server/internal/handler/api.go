package handler

import (
	"errors"
	"strconv"
	"strings"

	"pet-social-game/server/internal/middleware"
	"pet-social-game/server/internal/pkg/response"
	"pet-social-game/server/internal/service"

	"github.com/gin-gonic/gin"
)

type API struct {
	Svc *service.Services
}

func New(svc *service.Services) *API {
	return &API{Svc: svc}
}

func (a *API) Healthz(c *gin.Context) {
	response.OK(c, gin.H{"status": "up"})
}

func (a *API) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	resp, err := a.Svc.Login(c.Request.Context(), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, resp)
}

func (a *API) Me(c *gin.Context) {
	resp, err := a.Svc.Me(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, resp)
}

func (a *API) CreatePet(c *gin.Context) {
	var req service.CreatePetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	pet, err := a.Svc.CreatePet(middleware.UserID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, pet)
}

func (a *API) Feed(c *gin.Context) {
	var req service.FeedReq
	_ = c.ShouldBindJSON(&req)
	pet, err := a.Svc.Feed(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, pet)
}

func (a *API) Outfit(c *gin.Context) {
	var req service.OutfitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	pet, err := a.Svc.Outfit(middleware.UserID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, pet)
}

func (a *API) ListFriends(c *gin.Context) {
	list, err := a.Svc.ListFriends(middleware.UserID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, list)
}

func (a *API) RequestFriend(c *gin.Context) {
	var req service.FriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := a.Svc.RequestFriend(middleware.UserID(c), req); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (a *API) AcceptFriend(c *gin.Context) {
	var req service.FriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if err := a.Svc.AcceptFriend(middleware.UserID(c), req); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (a *API) VisitHome(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid userId")
		return
	}
	resp, err := a.Svc.VisitHome(middleware.UserID(c), hostID)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, resp)
}

func (a *API) VisitInteract(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid userId")
		return
	}
	var req service.InteractReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	pet, err := a.Svc.VisitInteract(middleware.UserID(c), hostID, req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, pet)
}

func (a *API) ChatHistory(c *gin.Context) {
	var q service.ChatHistoryQuery
	_ = c.ShouldBindQuery(&q)
	list, err := a.Svc.ChatHistory(middleware.UserID(c), q)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, list)
}

func writeErr(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case errors.Is(err, service.ErrBadRequest):
		response.BadRequest(c, trimErr(msg))
	case errors.Is(err, service.ErrUnauthorized):
		response.Unauthorized(c, trimErr(msg))
	case errors.Is(err, service.ErrForbidden):
		response.Forbidden(c, trimErr(msg))
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, trimErr(msg))
	case errors.Is(err, service.ErrConflict):
		response.Fail(c, 409, 409, trimErr(msg))
	default:
		response.ServerError(c, msg)
	}
}

func trimErr(msg string) string {
	for _, p := range []string{"bad request: ", "forbidden: ", "conflict: "} {
		if strings.HasPrefix(msg, p) {
			return strings.TrimPrefix(msg, p)
		}
	}
	return msg
}
