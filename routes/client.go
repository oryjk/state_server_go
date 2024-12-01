package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

type Status struct {
	ClientID string `json:"client_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

type ClientState struct {
	ID        uint      `gorm:"primaryKey"`
	ClientID  string    `gorm:"size:255;not null"`
	Status    string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"` // 使用 GORM 自动填充时间戳
}

// 显式指定表名
func (ClientState) TableName() string {
	return "ss_client_state_log" // 数据库中的实际表名
}

var (
	queue   = make(chan Status, 1000) // Buffered channel for the queue
	dbMutex sync.Mutex                // Mutex for database operations
)

// BatchInsert inserts a batch of statuses into the database
func BatchInsert(db *gorm.DB, statuses []Status) {
	if len(statuses) == 0 {
		return
	}
	clientStates := make([]ClientState, len(statuses))
	for i, s := range statuses {
		clientStates[i] = ClientState{
			ClientID: s.ClientID,
			Status:   s.Status,
		}
	}

	if err := db.Create(&clientStates).Error; err != nil {
		log.Fatalf("Failed to batch insert: %v", err)
	}

	log.Println("Batch insert successful!")

}

// ProcessQueue processes statuses from the queue in batches
func ProcessQueue(db *gorm.DB) {
	batchSize := 100
	ticker := time.NewTicker(500 * time.Millisecond) // Periodic batch processing

	for {
		select {
		case <-ticker.C:
			batch := make([]Status, 0, batchSize)

			for i := 0; i < batchSize; i++ {
				select {
				case status := <-queue:
					batch = append(batch, status)
				default:
					break
				}
			}

			if len(batch) > 0 {
				dbMutex.Lock()
				BatchInsert(db, batch)
				dbMutex.Unlock()

				log.Printf("Inserted %d records", len(batch))

			}
		}
	}
}

// ReceiveStatus handles incoming status requests
func ReceiveStatus(c *gin.Context) {
	var status Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	select {
	case queue <- status:
		c.JSON(http.StatusOK, gin.H{"message": "Status received"})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Queue is full"})
	}
}