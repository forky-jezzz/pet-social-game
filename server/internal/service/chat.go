package service

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"pet-social-game/server/internal/model"
)

type ChatHistoryQuery struct {
	Channel string `form:"channel"`
	PeerID  uint64 `form:"peerId"`
	Limit   int    `form:"limit"`
}

type ChatMsgView struct {
	ID       uint64 `json:"id"`
	Channel  string `json:"channel"`
	FromID   uint64 `json:"fromId"`
	FromName string `json:"fromName"`
	ToID     uint64 `json:"toId,omitempty"`
	Content  string `json:"content"`
	Ts       int64  `json:"ts"`
}

func (s *Services) ChatHistory(userID uint64, q ChatHistoryQuery) ([]ChatMsgView, error) {
	channel := q.Channel
	if channel == "" {
		channel = model.ChannelWorld
	}
	var msgs []model.ChatMessage
	var err error
	switch channel {
	case model.ChannelWorld:
		msgs, err = s.Store.ListWorldHistory(q.Limit)
	case model.ChannelFriend:
		if q.PeerID == 0 {
			return nil, fmt.Errorf("%w: peerId required", ErrBadRequest)
		}
		ok, e := s.Store.AreFriends(userID, q.PeerID)
		if e != nil {
			return nil, e
		}
		if !ok {
			return nil, ErrForbidden
		}
		msgs, err = s.Store.ListDMHistory(userID, q.PeerID, q.Limit)
	default:
		return nil, fmt.Errorf("%w: invalid channel", ErrBadRequest)
	}
	if err != nil {
		return nil, err
	}
	out := make([]ChatMsgView, 0, len(msgs))
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		name := ""
		if u, _ := s.Store.GetUser(m.FromID); u != nil {
			name = u.Nickname
		}
		v := ChatMsgView{
			ID: m.ID, Channel: m.Channel, FromID: m.FromID, FromName: name,
			Content: m.Content, Ts: m.CreatedAt.UnixMilli(),
		}
		if m.ToID != nil {
			v.ToID = *m.ToID
		}
		out = append(out, v)
	}
	return out, nil
}

type SendChatInput struct {
	Channel string
	ToID    uint64
	Content string
}

func (s *Services) SendChat(ctx context.Context, fromID uint64, in SendChatInput) (*ChatMsgView, string, error) {
	content := strings.TrimSpace(in.Content)
	if content == "" || utf8.RuneCountInString(content) > 200 {
		return nil, "", fmt.Errorf("%w: invalid content", ErrBadRequest)
	}
	channel := in.Channel
	if channel == "" {
		channel = model.ChannelWorld
	}

	pass, reason, err := s.Detector.DetectText(ctx, content)
	if err != nil {
		return nil, "", err
	}
	if !pass {
		if reason == "" {
			reason = "content_risk"
		}
		return nil, reason, nil
	}

	from, err := s.Store.GetUser(fromID)
	if err != nil || from == nil {
		return nil, "", ErrNotFound
	}

	msg := &model.ChatMessage{
		Channel: channel,
		FromID:  fromID,
		Content: content,
	}

	switch channel {
	case model.ChannelWorld:
		// ok
	case model.ChannelFriend:
		if in.ToID == 0 || in.ToID == fromID {
			return nil, "", fmt.Errorf("%w: invalid toId", ErrBadRequest)
		}
		ok, e := s.Store.AreFriends(fromID, in.ToID)
		if e != nil {
			return nil, "", e
		}
		if !ok {
			return nil, "", ErrForbidden
		}
		to := in.ToID
		msg.ToID = &to
	default:
		return nil, "", fmt.Errorf("%w: invalid channel", ErrBadRequest)
	}

	if err := s.Store.CreateChatMessage(msg); err != nil {
		return nil, "", err
	}
	view := &ChatMsgView{
		ID: msg.ID, Channel: msg.Channel, FromID: fromID, FromName: from.Nickname,
		Content: content, Ts: msg.CreatedAt.UnixMilli(),
	}
	if msg.ToID != nil {
		view.ToID = *msg.ToID
	}
	return view, "", nil
}
