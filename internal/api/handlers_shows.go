package api

import (
	"log/slog"
	"net/http"
	"uuid"

	"github.com/bharsadia-parth/seatres/internal/store"
	"github.com/gin-gonic/gin"
)

type createShowReq struct {
	Name         string   `json:"name" binding:"required"`
	Seats        []string `json:"seats" binding:"required,min=1,max=200000"`
	PricePaise   int64    `json:"price_paise" binding:"required,gt=0"`
	PerUserLimit *int     `json:"per_user_limit"` // optional, default 4
}

func CreateShow(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createShowReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "data": err.Error()})
			return
		}

		limit := 4
		if req.PerUserLimit != nil {
			if *req.PerUserLimit < 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "per_user_limit_must_be_positive"})
				return
			}
			limit = *req.PerUserLimit
		}
		show, err := st.CreateShow(c.Request.Context(), req.Name, req.Seats, req.PricePaise, limit)
		if err != nil {
			slog.Error("get show failed", "request_id", c.GetString("request_id"), "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
			return
		}
		c.JSON(http.StatusCreated, show)
	}
}

func GetShow(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		show, err := st.GetShow(c.Request.Context(), id)
		if err != nil {
			slog.Error("get show failed", "request_id", c.GetString("request_id"), "err", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusOK, show)
	}
}
