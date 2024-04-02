package handler

import (
	"flotify/internal/repository"
	"flotify/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(dbpool *pgxpool.Pool) *gin.Engine {
	router := gin.New()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	artist_repo := repository.NewPostgresArtistRepository(dbpool)
	artist_handler := NewArtistHandler(artist_repo)
	artist_subrouter := router.Group("/artists")
	{
		artist_subrouter.POST("/", artist_handler.CreateArtist)
		artist_subrouter.GET("/:id", artist_handler.GetInfoArtistByID)
		artist_subrouter.GET("/:id/tracks", artist_handler.GetArtistTracksByID)
		artist_subrouter.PUT("/", artist_handler.UpdateArtist)
		artist_subrouter.DELETE("/:id", artist_handler.DeleteArtist)
		artist_subrouter.GET("/", artist_handler.GetArtistWithFilter)
	}

	user_repo := repository.NewPostgresUserRepository(dbpool)
	user_handler := NewUserHandler(user_repo)
	user_subrouter := router.Group("/users")
	{
		user_subrouter.POST("/register", user_handler.CreateUser)
		user_subrouter.POST("/login", user_handler.LoginUser)
		user_subrouter.Use(middleware.AuthRequest())
		user_subrouter.GET("/:id", user_handler.ViewInformation)
	}

	return router
}
