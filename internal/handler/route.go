package handler

import (
	"flotify/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(dbpool *pgxpool.Pool) *gin.Engine {
	router := gin.New()

	track_repo := repository.NewPostgresTrackRepository(dbpool)
	track_handler := NewTrackHandler(track_repo)
	track_subrouter := router.Group("/tracks")
	{
		track_subrouter.POST("/", track_handler.CreateTrack)
		track_subrouter.GET("/:id", track_handler.GetTrackByID)
		track_subrouter.PUT("/", track_handler.UpdateTrack)
		track_subrouter.DELETE("/:id", track_handler.DeleteTrack)
		track_subrouter.GET("/", track_handler.GetTrackWithFilter)
	}

	return router
}
