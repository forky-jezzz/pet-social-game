package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pet-social-game/server/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Store struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewStore(db *gorm.DB, rdb *redis.Client) *Store {
	return &Store{DB: db, RDB: rdb}
}

func (s *Store) FindUserByPlatformOpenID(platform, openID string) (*model.User, error) {
	var u model.User
	err := s.DB.Where("platform = ? AND open_id = ?", platform, openID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (s *Store) CreateUser(u *model.User) error {
	return s.DB.Create(u).Error
}

func (s *Store) GetUser(id uint64) (*model.User, error) {
	var u model.User
	err := s.DB.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, err
}

func (s *Store) GetPetByUser(userID uint64) (*model.Pet, error) {
	var p model.Pet
	err := s.DB.Where("user_id = ?", userID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (s *Store) CreatePet(p *model.Pet) error {
	return s.DB.Create(p).Error
}

func (s *Store) UpdatePet(p *model.Pet) error {
	return s.DB.Save(p).Error
}

func (s *Store) WithFeedLock(ctx context.Context, userID uint64, fn func() error) error {
	key := fmt.Sprintf("lock:feed:%d", userID)
	ok, err := s.RDB.SetNX(ctx, key, "1", 3*time.Second).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("feed busy")
	}
	defer s.RDB.Del(ctx, key)
	return fn()
}

func (s *Store) SetOnline(ctx context.Context, userID uint64, ttl time.Duration) error {
	return s.RDB.Set(ctx, fmt.Sprintf("online:%d", userID), "1", ttl).Err()
}

func (s *Store) IsOnline(ctx context.Context, userID uint64) (bool, error) {
	n, err := s.RDB.Exists(ctx, fmt.Sprintf("online:%d", userID)).Result()
	return n > 0, err
}

func (s *Store) ListFriends(userID uint64) ([]model.Friendship, error) {
	var list []model.Friendship
	err := s.DB.Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, model.FriendAccepted).
		Find(&list).Error
	return list, err
}

func (s *Store) ListPendingIncoming(userID uint64) ([]model.Friendship, error) {
	var list []model.Friendship
	err := s.DB.Where("friend_id = ? AND status = ?", userID, model.FriendPending).Find(&list).Error
	return list, err
}

func (s *Store) GetFriendship(a, b uint64) (*model.Friendship, error) {
	var f model.Friendship
	err := s.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", a, b, b, a).
		First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &f, err
}

func (s *Store) CreateFriendship(f *model.Friendship) error {
	return s.DB.Create(f).Error
}

func (s *Store) UpdateFriendship(f *model.Friendship) error {
	return s.DB.Save(f).Error
}

func (s *Store) AreFriends(a, b uint64) (bool, error) {
	f, err := s.GetFriendship(a, b)
	if err != nil || f == nil {
		return false, err
	}
	return f.Status == model.FriendAccepted, nil
}

func (s *Store) CreateVisitLog(v *model.VisitLog) error {
	return s.DB.Create(v).Error
}

func (s *Store) CreateChatMessage(m *model.ChatMessage) error {
	return s.DB.Create(m).Error
}

func (s *Store) ListWorldHistory(limit int) ([]model.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var list []model.ChatMessage
	err := s.DB.Where("channel = ?", model.ChannelWorld).
		Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (s *Store) ListDMHistory(a, b uint64, limit int) ([]model.ChatMessage, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var list []model.ChatMessage
	err := s.DB.Where(
		"channel = ? AND ((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?))",
		model.ChannelFriend, a, b, b, a,
	).Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (s *Store) GetItemByCode(code string) (*model.Item, error) {
	var it model.Item
	err := s.DB.Where("code = ?", code).First(&it).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &it, err
}

func (s *Store) ListOutfitItems() ([]model.Item, error) {
	var list []model.Item
	err := s.DB.Where("category = ?", "outfit").Find(&list).Error
	return list, err
}
