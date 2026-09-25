package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pet-social-game/server/internal/config"
	"pet-social-game/server/internal/model"
	"pet-social-game/server/internal/pkg/auth"
	"pet-social-game/server/internal/repo"
	"pet-social-game/server/internal/security"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
)

type Services struct {
	Cfg      *config.Config
	Store    *repo.Store
	Detector security.Detector
}

func New(cfg *config.Config, store *repo.Store, detector security.Detector) *Services {
	return &Services{Cfg: cfg, Store: store, Detector: detector}
}

type LoginReq struct {
	Platform string `json:"platform"`
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type LoginResp struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
	Pet   *model.Pet  `json:"pet"`
}

func (s *Services) Login(ctx context.Context, req LoginReq) (*LoginResp, error) {
	platform := strings.ToLower(strings.TrimSpace(req.Platform))
	code := strings.TrimSpace(req.Code)
	if platform == "" || code == "" {
		return nil, fmt.Errorf("%w: platform and code required", ErrBadRequest)
	}
	if platform == "dev" && !s.Cfg.Auth.AllowDevLogin {
		return nil, fmt.Errorf("%w: dev login disabled", ErrForbidden)
	}

	openID, err := s.resolveOpenID(platform, code)
	if err != nil {
		return nil, err
	}

	user, err := s.Store.FindUserByPlatformOpenID(platform, openID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		nick := req.Nickname
		if nick == "" {
			nick = "旅行者" + openID[len(openID)-4:]
		}
		user = &model.User{
			OpenID:   openID,
			Platform: platform,
			Nickname: nick,
			Avatar:   req.Avatar,
		}
		if err := s.Store.CreateUser(user); err != nil {
			return nil, err
		}
	}

	token, err := auth.IssueToken(s.Cfg.JWT.Secret, user.ID, user.Platform, s.Cfg.JWT.ExpireHours)
	if err != nil {
		return nil, err
	}
	pet, _ := s.Store.GetPetByUser(user.ID)
	_ = s.Store.SetOnline(ctx, user.ID, 5*time.Minute)
	return &LoginResp{Token: token, User: user, Pet: pet}, nil
}

func (s *Services) resolveOpenID(platform, code string) (string, error) {
	switch platform {
	case "dev":
		return "dev_" + code, nil
	case "alipay", "wechat":
		// Production: exchange code with platform OAuth.
		// Dev-friendly: treat code as open_id suffix when credentials absent.
		return platform + "_" + code, nil
	default:
		return "", fmt.Errorf("%w: unsupported platform", ErrBadRequest)
	}
}

type MeResp struct {
	User   *model.User `json:"user"`
	Pet    *model.Pet  `json:"pet"`
	Online bool        `json:"online"`
}

func (s *Services) Me(ctx context.Context, userID uint64) (*MeResp, error) {
	user, err := s.Store.GetUser(userID)
	if err != nil || user == nil {
		return nil, ErrNotFound
	}
	pet, err := s.Store.GetPetByUser(userID)
	if err != nil {
		return nil, err
	}
	online, _ := s.Store.IsOnline(ctx, userID)
	_ = s.Store.SetOnline(ctx, userID, 5*time.Minute)
	return &MeResp{User: user, Pet: pet, Online: online}, nil
}

type CreatePetReq struct {
	Name string `json:"name"`
}

func (s *Services) CreatePet(userID uint64, req CreatePetReq) (*model.Pet, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len([]rune(name)) > 16 {
		return nil, fmt.Errorf("%w: invalid pet name", ErrBadRequest)
	}
	exist, err := s.Store.GetPetByUser(userID)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, fmt.Errorf("%w: pet already exists", ErrConflict)
	}
	outfit := map[string]string{"hat": "hat_basic", "clothes": "clothes_basic"}
	b, _ := json.Marshal(outfit)
	pet := &model.Pet{
		UserID:     userID,
		Name:       name,
		Level:      1,
		Hunger:     40,
		Mood:       60,
		OutfitJSON: string(b),
	}
	if err := s.Store.CreatePet(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

type FeedReq struct {
	FoodCode string `json:"foodCode"`
}

func (s *Services) Feed(ctx context.Context, userID uint64, req FeedReq) (*model.Pet, error) {
	foodCode := req.FoodCode
	if foodCode == "" {
		foodCode = "food_kibble"
	}
	item, err := s.Store.GetItemByCode(foodCode)
	if err != nil {
		return nil, err
	}
	if item == nil || item.Category != "food" {
		return nil, fmt.Errorf("%w: invalid food", ErrBadRequest)
	}

	hungerAdd, moodAdd := 20, 5
	if item.MetaJSON != nil {
		var meta map[string]int
		_ = json.Unmarshal([]byte(*item.MetaJSON), &meta)
		if v, ok := meta["hunger"]; ok {
			hungerAdd = v
		}
		if v, ok := meta["mood"]; ok {
			moodAdd = v
		}
	}

	var pet *model.Pet
	err = s.Store.WithFeedLock(ctx, userID, func() error {
		p, e := s.Store.GetPetByUser(userID)
		if e != nil {
			return e
		}
		if p == nil {
			return ErrNotFound
		}
		p.Hunger = clamp(p.Hunger+hungerAdd, 0, 100)
		p.Mood = clamp(p.Mood+moodAdd, 0, 100)
		if p.Hunger >= 80 && p.Mood >= 80 {
			p.Level++
			p.Hunger = 50
			p.Mood = 50
		}
		if e := s.Store.UpdatePet(p); e != nil {
			return e
		}
		pet = p
		return nil
	})
	if err != nil {
		if err.Error() == "feed busy" {
			return nil, fmt.Errorf("%w: feed busy", ErrConflict)
		}
		return nil, err
	}
	return pet, nil
}

type OutfitReq struct {
	Hat     string `json:"hat"`
	Clothes string `json:"clothes"`
}

func (s *Services) Outfit(userID uint64, req OutfitReq) (*model.Pet, error) {
	pet, err := s.Store.GetPetByUser(userID)
	if err != nil || pet == nil {
		return nil, ErrNotFound
	}
	outfit := map[string]string{}
	_ = json.Unmarshal([]byte(pet.OutfitJSON), &outfit)
	if req.Hat != "" {
		if it, _ := s.Store.GetItemByCode(req.Hat); it == nil || it.Category != "outfit" {
			return nil, fmt.Errorf("%w: invalid hat", ErrBadRequest)
		}
		outfit["hat"] = req.Hat
	}
	if req.Clothes != "" {
		if it, _ := s.Store.GetItemByCode(req.Clothes); it == nil || it.Category != "outfit" {
			return nil, fmt.Errorf("%w: invalid clothes", ErrBadRequest)
		}
		outfit["clothes"] = req.Clothes
	}
	b, _ := json.Marshal(outfit)
	pet.OutfitJSON = string(b)
	if err := s.Store.UpdatePet(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
