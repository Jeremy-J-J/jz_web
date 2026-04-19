package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"jz_web/utils"
)

type StatsResponse struct {
	Total   int64 `json:"total"`
	Today   int64 `json:"today"`
}

func GetStats(c *gin.Context) {
	var total int64
	err := utils.DB.QueryRow("SELECT COALESCE(SUM(count), 0) FROM stats").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total stats"})
		return
	}

	today := time.Now().Format("2006-01-02")
	var todayCount int64
	err = utils.DB.QueryRow("SELECT COALESCE(count, 0) FROM stats WHERE date = ?", today).Scan(&todayCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get today's stats"})
		return
	}

	c.JSON(http.StatusOK, StatsResponse{
		Total: total,
		Today: todayCount,
	})
}

func RecordVisit(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var existingCount int
	err := utils.DB.QueryRow("SELECT count FROM stats WHERE date = ?", today).Scan(&existingCount)
	if err != nil {
		utils.DB.Exec("INSERT INTO stats (date, count) VALUES (?, 1)", today)
	} else {
		utils.DB.Exec("UPDATE stats SET count = count + 1 WHERE date = ?", today)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visit recorded"})
}
