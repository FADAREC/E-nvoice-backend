package models

import "time"

type MutationLog struct {
    ID              int64           `json:"id"`
    UserID          string          `json:"user_id"`
    EntityType      string          `json:"entity_type"`
    EntityID        string          `json:"entity_id"`
    Operation       string          `json:"operation"`
    Payload         interface{}     `json:"payload"`
    DeviceID        string          `json:"device_id"`
    ClientTimestamp time.Time       `json:"client_timestamp"`
    ServerTimestamp time.Time       `json:"server_timestamp"`
}

type SyncPushRequest struct {
    DeviceID  string        `json:"device_id" binding:"required"`
    Mutations []SyncMutation `json:"mutations" binding:"required"`
}

type SyncMutation struct {
    EntityType      string      `json:"entity_type" binding:"required"`
    EntityID        string      `json:"entity_id" binding:"required"`
    Operation       string      `json:"operation" binding:"required"`
    Payload         interface{} `json:"payload" binding:"required"`
    ClientTimestamp time.Time   `json:"client_timestamp" binding:"required"`
}

type SyncPullRequest struct {
    Since    *time.Time `json:"since"`
    DeviceID string     `json:"device_id" binding:"required"`
}

type SyncPullResponse struct {
    Mutations      []MutationLog `json:"mutations"`
    ServerTime     time.Time     `json:"server_time"`
    HasMore        bool          `json:"has_more"`
}