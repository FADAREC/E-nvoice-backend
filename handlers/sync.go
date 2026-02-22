package handlers

import (
    // "database/sql"
    "encoding/json"
    "net/http"
    "time"
    "invoice-backend/database"
    "invoice-backend/models"
    "github.com/gin-gonic/gin"
)

type SyncHandler struct{}

func (h *SyncHandler) Push(c *gin.Context) {
    userID := c.GetString("user_id")

    var req models.SyncPushRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tx, err := database.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }
    defer tx.Rollback()

    for _, mutation := range req.Mutations {
        payloadJSON, err := json.Marshal(mutation.Payload)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
            return
        }

        _, err = tx.Exec(`
            INSERT INTO mutation_log 
            (user_id, entity_type, entity_id, operation, payload, device_id, client_timestamp)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, userID, mutation.EntityType, mutation.EntityID, mutation.Operation, payloadJSON, req.DeviceID, mutation.ClientTimestamp)

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert mutation"})
            return
        }
    }

    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "synced": len(req.Mutations),
        "timestamp": time.Now(),
    })
}

func (h *SyncHandler) Pull(c *gin.Context) {
    userID := c.GetString("user_id")

    var req models.SyncPullRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `
        SELECT id, user_id, entity_type, entity_id, operation, payload, device_id, client_timestamp, server_timestamp
        FROM mutation_log
        WHERE user_id = $1 AND device_id != $2
    `
    args := []interface{}{userID, req.DeviceID}

    if req.Since != nil {
        query += " AND server_timestamp > $3"
        args = append(args, req.Since)
    }

    query += " ORDER BY server_timestamp ASC LIMIT 100"

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mutations"})
        return
    }
    defer rows.Close()

    var mutations []models.MutationLog
    for rows.Next() {
        var m models.MutationLog
        var payloadJSON []byte

        err := rows.Scan(
            &m.ID, &m.UserID, &m.EntityType, &m.EntityID, &m.Operation,
            &payloadJSON, &m.DeviceID, &m.ClientTimestamp, &m.ServerTimestamp,
        )
        if err != nil {
            continue
        }

        json.Unmarshal(payloadJSON, &m.Payload)
        mutations = append(mutations, m)
    }

    c.JSON(http.StatusOK, models.SyncPullResponse{
        Mutations:  mutations,
        ServerTime: time.Now(),
        HasMore:    len(mutations) == 100,
    })
}