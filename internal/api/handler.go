package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NeverovDS/playlist/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Player
}

type InputSong struct {
	ID       int
	Name     string `json:"name" binding:"required"`
	Duration int    `json:"duration" binding:"required"`
}

func NewHandler(services *service.Player) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	playlist_service := router.Group("/playlist")
	{
		playlist_service.GET("/play", func(c *gin.Context) {
			err := h.services.Playlist.Play()
			if err != nil {
				newResponse(c, http.StatusBadRequest, err.Error())
			} else {
				c.String(200, "Плейлист запущен!")
			}
		})

		playlist_service.GET("/pause", func(c *gin.Context) {
			err := h.services.Playlist.Pause()
			if err != nil {
				newResponse(c, http.StatusBadRequest, err.Error())
			} else {
				c.String(200, "Воспроизведение плейлиста приостановлено")
			}
		})

		playlist_service.GET("/next", func(c *gin.Context) {
			err := h.services.Playlist.Next()
			if err != nil {
				newResponse(c, http.StatusBadRequest, err.Error())
			} else {
				c.String(200, "Следующий трек включен")
			}
		})

		playlist_service.GET("/prev", func(c *gin.Context) {
			err := h.services.Playlist.Prev()
			if err != nil {
				newResponse(c, http.StatusBadRequest, err.Error())
			} else {
				c.String(200, "Предыдущий трек включен")
			}
		})

		playlist_service.POST("/create/song", h.createSong)

		playlist_service.GET("/songs", h.songsList)

		playlist_service.GET("/song/:id", h.showSong)

		playlist_service.PUT("/song/:id", h.updateSong)

		playlist_service.DELETE("/song/:id", h.deleteSong)
	}
	return router
}

func (h *Handler) createSong(c *gin.Context) {
	var inp InputSong
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Playlist.WriteSongToJson(&service.Song{
		Name:     inp.Name,
		Duration: time.Duration(inp.Duration),
	})
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	}
	h.services.Playlist.CreateSong(&service.Song{
		Name:     inp.Name,
		Duration: time.Duration(inp.Duration) * time.Second,
	})
	c.String(200, "Песня %s добавлена в плейлист", inp.Name)
}
func (h *Handler) songsList(c *gin.Context) {
	err, jsonData := service.UnmarshalJson()
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	} else {
		c.JSON(200, jsonData)
	}

}
func (h *Handler) showSong(c *gin.Context) {
	idInt, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	}
	song, err := h.services.Playlist.ShowSong(&service.Song{ID: idInt})
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	} else {
		c.JSON(200, song)
	}
}

func (h *Handler) updateSong(c *gin.Context) {

	idInt, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	}
	var inp InputSong

	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	song := h.services.Playlist.UpdateSong(&service.Song{ID: idInt, Name: inp.Name, Duration: time.Duration(inp.Duration)})

	c.JSON(200, song)
}
func (h *Handler) deleteSong(c *gin.Context) {

	idInt, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		panic(err)
	}

	err = h.services.Playlist.DeleteSong(&service.Song{ID: idInt})
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
	} else {
		c.String(200, "Song deleted")
	}
}
