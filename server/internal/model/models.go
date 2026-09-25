package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	OpenID    string    `gorm:"column:open_id;size:128;not null" json:"openId"`
	Platform  string    `gorm:"size:32;not null" json:"platform"`
	Nickname  string    `gorm:"size:64;not null" json:"nickname"`
	Avatar    string    `gorm:"size:512;not null" json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string { return "users" }

type Pet struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `gorm:"column:user_id;uniqueIndex;not null" json:"userId"`
	Name       string    `gorm:"size:64;not null" json:"name"`
	Level      int       `gorm:"not null;default:1" json:"level"`
	Hunger     int       `gorm:"not null;default:50" json:"hunger"`
	Mood       int       `gorm:"not null;default:50" json:"mood"`
	OutfitJSON string    `gorm:"column:outfit_json;type:json;not null" json:"outfitJson"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Pet) TableName() string { return "pets" }

type Item struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Category  string    `gorm:"size:32;not null" json:"category"`
	MetaJSON  *string   `gorm:"column:meta_json;type:json" json:"metaJson,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Item) TableName() string { return "items" }

type UserItem struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"column:user_id;not null" json:"userId"`
	ItemID    uint64    `gorm:"column:item_id;not null" json:"itemId"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserItem) TableName() string { return "user_items" }

const (
	FriendPending  = "pending"
	FriendAccepted = "accepted"
)

type Friendship struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"column:user_id;not null" json:"userId"`
	FriendID  uint64    `gorm:"column:friend_id;not null" json:"friendId"`
	Status    string    `gorm:"size:16;not null" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Friendship) TableName() string { return "friendships" }

const (
	ChannelWorld  = "world"
	ChannelFriend = "friend"
)

type ChatMessage struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Channel   string    `gorm:"size:16;not null" json:"channel"`
	FromID    uint64    `gorm:"column:from_id;not null" json:"fromId"`
	ToID      *uint64   `gorm:"column:to_id" json:"toId,omitempty"`
	Content   string    `gorm:"size:512;not null" json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ChatMessage) TableName() string { return "chat_messages" }

const (
	VisitView = "view"
	VisitPet  = "pet"
	VisitGift = "gift"
)

type VisitLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	VisitorID uint64    `gorm:"column:visitor_id;not null" json:"visitorId"`
	HostID    uint64    `gorm:"column:host_id;not null" json:"hostId"`
	Action    string    `gorm:"size:32;not null" json:"action"`
	CreatedAt time.Time `json:"createdAt"`
}

func (VisitLog) TableName() string { return "visit_logs" }
