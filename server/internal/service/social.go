package service

import (
	"context"
	"fmt"

	"pet-social-game/server/internal/model"
)

type FriendView struct {
	UserID   uint64 `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
	Online   bool   `json:"online"`
}

func (s *Services) ListFriends(userID uint64) ([]FriendView, error) {
	list, err := s.Store.ListFriends(userID)
	if err != nil {
		return nil, err
	}
	pending, err := s.Store.ListPendingIncoming(userID)
	if err != nil {
		return nil, err
	}
	out := make([]FriendView, 0, len(list)+len(pending))
	for _, f := range list {
		other := f.FriendID
		if f.FriendID == userID {
			other = f.UserID
		}
		u, err := s.Store.GetUser(other)
		if err != nil || u == nil {
			continue
		}
		online, _ := s.Store.IsOnline(context.Background(), other)
		out = append(out, FriendView{
			UserID: u.ID, Nickname: u.Nickname, Avatar: u.Avatar, Status: model.FriendAccepted, Online: online,
		})
	}
	for _, f := range pending {
		u, err := s.Store.GetUser(f.UserID)
		if err != nil || u == nil {
			continue
		}
		out = append(out, FriendView{
			UserID: u.ID, Nickname: u.Nickname, Avatar: u.Avatar, Status: model.FriendPending, Online: false,
		})
	}
	return out, nil
}

type FriendReq struct {
	FriendID uint64 `json:"friendId"`
}

func (s *Services) RequestFriend(userID uint64, req FriendReq) error {
	if req.FriendID == 0 || req.FriendID == userID {
		return fmt.Errorf("%w: invalid friendId", ErrBadRequest)
	}
	target, err := s.Store.GetUser(req.FriendID)
	if err != nil || target == nil {
		return ErrNotFound
	}
	exist, err := s.Store.GetFriendship(userID, req.FriendID)
	if err != nil {
		return err
	}
	if exist != nil {
		if exist.Status == model.FriendAccepted {
			return fmt.Errorf("%w: already friends", ErrConflict)
		}
		return fmt.Errorf("%w: request already exists", ErrConflict)
	}
	return s.Store.CreateFriendship(&model.Friendship{
		UserID:   userID,
		FriendID: req.FriendID,
		Status:   model.FriendPending,
	})
}

func (s *Services) AcceptFriend(userID uint64, req FriendReq) error {
	if req.FriendID == 0 {
		return fmt.Errorf("%w: invalid friendId", ErrBadRequest)
	}
	f, err := s.Store.GetFriendship(userID, req.FriendID)
	if err != nil {
		return err
	}
	if f == nil || f.Status != model.FriendPending {
		return ErrNotFound
	}
	// Only the invitee (friend_id) can accept.
	if f.FriendID != userID {
		return ErrForbidden
	}
	f.Status = model.FriendAccepted
	return s.Store.UpdateFriendship(f)
}

type VisitHomeResp struct {
	Host   *model.User `json:"host"`
	Pet    *model.Pet  `json:"pet"`
	Online bool        `json:"online"`
}

func (s *Services) VisitHome(visitorID, hostID uint64) (*VisitHomeResp, error) {
	if hostID == 0 {
		return nil, fmt.Errorf("%w: invalid host", ErrBadRequest)
	}
	host, err := s.Store.GetUser(hostID)
	if err != nil || host == nil {
		return nil, ErrNotFound
	}
	pet, err := s.Store.GetPetByUser(hostID)
	if err != nil {
		return nil, err
	}
	online, _ := s.Store.IsOnline(context.Background(), hostID)
	_ = s.Store.CreateVisitLog(&model.VisitLog{
		VisitorID: visitorID,
		HostID:    hostID,
		Action:    model.VisitView,
	})
	return &VisitHomeResp{Host: host, Pet: pet, Online: online}, nil
}

type InteractReq struct {
	Action string `json:"action"` // pet|gift
}

func (s *Services) VisitInteract(visitorID, hostID uint64, req InteractReq) (*model.Pet, error) {
	if hostID == 0 || hostID == visitorID {
		return nil, fmt.Errorf("%w: invalid host", ErrBadRequest)
	}
	action := req.Action
	if action != model.VisitPet && action != model.VisitGift {
		return nil, fmt.Errorf("%w: invalid action", ErrBadRequest)
	}
	host, err := s.Store.GetUser(hostID)
	if err != nil || host == nil {
		return nil, ErrNotFound
	}
	pet, err := s.Store.GetPetByUser(hostID)
	if err != nil || pet == nil {
		return nil, ErrNotFound
	}
	if action == model.VisitPet {
		pet.Mood = clamp(pet.Mood+8, 0, 100)
	} else {
		pet.Hunger = clamp(pet.Hunger+10, 0, 100)
		pet.Mood = clamp(pet.Mood+5, 0, 100)
	}
	if err := s.Store.UpdatePet(pet); err != nil {
		return nil, err
	}
	_ = s.Store.CreateVisitLog(&model.VisitLog{
		VisitorID: visitorID,
		HostID:    hostID,
		Action:    action,
	})
	return pet, nil
}
