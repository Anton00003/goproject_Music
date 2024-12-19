package api

import (
	"context"
	"errors"
	"goproject_Music/datastruct"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	docs "goproject_Music/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type updateReq struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	GroupId int    `json:"group"`
	Date    string `json:"date"`
	Text    string `json:"text"`
	Link    string `json:"link"`
}

type addReq struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}
type deleteReq struct {
	Id int `json:"id"`
}

type serv interface {
	GetAllTextMusicById(ctx context.Context, id int) (string, error)
	GetPaginTextMusicById(ctx context.Context, id, nOnPage int, nPage int) ([]string, error)
	AddMusic(ctx context.Context, song string, group string) (*datastruct.Music, error)
	GetMusicByFilter(ctx context.Context, filter *datastruct.Music, nOnPage int, nPage int) ([]datastruct.Music, error)
	UpdateMusicById(ctx context.Context, m *datastruct.Music) error
	DeleteMusicById(ctx context.Context, id int) error
	GetGroupId(ctx context.Context, group string) (int, error)
	GetList(ctx context.Context) ([]datastruct.MusicListItem, error)
	GetSongFromClient(name, group string) (*datastruct.SongDetail, error)
}

type api struct {
	Serv serv
}

func NewApi(s serv) *api {
	a := &api{Serv: s}
	return a
}
func (a *api) Run(host string) {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.GET("/music/text", a.getAllTextMusicById)
	r.GET("/music/text/couplet", a.getPaginTextMusicById)
	r.GET("/music", a.getMusicByFilter)
	r.POST("/music", a.addMusic)
	r.PATCH("/music", a.updateMusicById)
	r.DELETE("/music", a.deleteMusicById)
	r.GET("/music/list", a.getList)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(host)
}

// AllTextMusic godoc
// @Summary Get AllTextMusic
// @Schemes
// @Description do AllTextMusic
// @Tags AllTextMusic
// @Accept json
// @Produce json
// @Param        id   query      int  false  "Id"
// @Success 200 {string} TextSong
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /music/text [get]
func (a *api) getAllTextMusicById(g *gin.Context) {
	ctx := g.Request.Context()
	idS := g.Request.URL.Query().Get("id")
	id, err := strconv.Atoi(idS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number id is not numerical")
		log.Error("Error number id parameter")
		return
	}

	if id <= 0 {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	TextSong, err := a.Serv.GetAllTextMusicById(ctx, id)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadId) { // ErrBadNameGroup заменить на ErrBadId
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, TextSong)
}

// PaginTextMusic godoc
// @Summary Get PaginTextMusic
// @Schemes
// @Description do PaginTextMusic
// @Tags PaginTextMusic
// @Accept json
// @Produce json
// @Param          id     query      int     false  "Id"
// @Param        nOnPage  query      int     false  "nOnPage"
// @Param        nPage    query      int     false  "nPage"
// @Success 200 {string} TextSong
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /music/text/couplet [get]
func (a *api) getPaginTextMusicById(g *gin.Context) {
	ctx := g.Request.Context()
	idS := g.Request.URL.Query().Get("id")
	nOnPageS := g.Request.URL.Query().Get("nOnPage")
	nPageS := g.Request.URL.Query().Get("nPage")

	id, err := strconv.Atoi(idS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number id is not numerical")
		log.Error("Error number id parameter")
		return
	}
	nOnPage, err := strconv.Atoi(nOnPageS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number records on page is not numerical")
		log.Error("Error number records on page parameter")
		return
	}
	nPage, err := strconv.Atoi(nPageS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number of page is not numerical")
		log.Error("Error number of page parameter")
		return
	}

	if id <= 0 {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	if nOnPage <= 0 {
		g.JSON(http.StatusBadRequest, "number records on page can not be <= 0")
		log.Error("Error Music parameter")
		return
	}
	if nPage <= 0 {
		g.JSON(http.StatusBadRequest, "number of page can not be <= 0")
		log.Error("Error Music parameter")
		return
	}

	coupOnPage, err := a.Serv.GetPaginTextMusicById(ctx, id, nOnPage, nPage)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadNumPage) || errors.Is(err, datastruct.ErrBadId) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, coupOnPage)
}

// FilterMusic godoc
// @Summary Get FilterMusic
// @Schemes
// @Description do FilterMusic
// @Tags FilterMusic
// @Accept json
// @Produce json
// @Param        id       query      int     false  "Id"
// @Param        name     query      string  false  "Name"
// @Param        group    query      int     false  "Group"
// @Param        date     query      string  false  "Date"
// @Param        text     query      string  false  "Text"
// @Param        link     query      string  false  "Link"
// @Param        nOnPage  query      int     false  "nOnPage"
// @Param        nPage    query      int     false  "nPage"
// @Success 200 {array} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /music [get]
func (a *api) getMusicByFilter(g *gin.Context) {
	ctx := g.Request.Context()
	var err error
	log.Debug("Api: Filter Run")
	filter := &datastruct.Music{}

	musicId := g.Request.URL.Query().Get("id")
	filter.Name = g.Request.URL.Query().Get("name")
	groupId := g.Request.URL.Query().Get("group")
	date := g.Request.URL.Query().Get("date")
	filter.Text = g.Request.URL.Query().Get("text")
	filter.Link = g.Request.URL.Query().Get("link")
	nOnPageS := g.Request.URL.Query().Get("nOnPage")
	nPageS := g.Request.URL.Query().Get("nPage")

	if musicId != "" {
		filter.Id, err = strconv.Atoi(musicId)
		if err != nil {
			g.JSON(http.StatusBadRequest, "number musicId is not numerical")
			log.Error("Error number musicId parameter")
			return
		}
	}
	if groupId != "" {
		filter.GroupId, err = strconv.Atoi(groupId)
		if err != nil {
			g.JSON(http.StatusBadRequest, "number groupId is not numerical")
			log.Error("Error number groupId parameter")
			return
		}
	}
	if date != "" {
		filter.Date, err = time.Parse("02.01.2006", date)
		if err != nil {
			g.JSON(http.StatusBadRequest, "input bad date")
			log.Error("Error, input bad date")
			return
		}
	}
	nOnPage, err := strconv.Atoi(nOnPageS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number records on page is not numerical")
		log.Error("Error number records on page parameter")
		return
	}
	nPage, err := strconv.Atoi(nPageS)
	if err != nil {
		g.JSON(http.StatusBadRequest, "number of page is not numerical")
		log.Error("Error number of page parameter")
		return
	}
	if filter.Id <= 0 && filter.Name == "" && filter.GroupId <= 0 && filter.Date.IsZero() && filter.Text == "" && filter.Link == "" {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	if nOnPage <= 0 {
		g.JSON(http.StatusBadRequest, "number records on page can not be <= 0")
		log.Error("Error Music parameter")
		return
	}
	if nPage <= 0 {
		g.JSON(http.StatusBadRequest, "page can not be <= 0")
		log.Error("Error Music parameter")
		return
	}

	songs, err := a.Serv.GetMusicByFilter(ctx, filter, nOnPage, nPage)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadFilter) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, songs)
}

// AddMusic godoc
// @Summary Add Music
// @Schemes
// @Description do Music
// @Tags Music
// @Accept json
// @Produce json
// @Param        p   body      addReq  true  "Music"
// @Success 200 {object} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 500 {string} ServerError
// @Router /music [post]
func (a *api) addMusic(g *gin.Context) {
	ctx := g.Request.Context()
	m := &datastruct.Music{}
	p := addReq{}
	err := g.ShouldBindJSON(&p)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body: ", err.Error())
		return
	}

	if p.Name == "" {
		g.JSON(http.StatusBadRequest, "name is required")
		log.Error("Error, name is required")
		return
	}
	if p.Group == "" {
		g.JSON(http.StatusBadRequest, "group is required")
		log.Error("Error, group is required")
		return
	}

	//	var m *datastruct.Music
	m, err = a.Serv.AddMusic(ctx, p.Name, p.Group)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Error, Serv.AddMusic: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, *m)
}

// UpdateMusic godoc
// @Summary Update Music
// @Schemes
// @Description do Update Music
// @Tags UpdateMusic
// @Accept json
// @Produce json
// @Param        m   body      updateReq  true  "Music"
// @Success 200 {object} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 500 {string} ServerError
// @Router /music [patch]
func (a *api) updateMusicById(g *gin.Context) {
	ctx := g.Request.Context()
	u := &updateReq{}

	err := g.ShouldBindJSON(u)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body", err.Error())
		return
	}
	if u.Id <= 0 {
		g.JSON(http.StatusBadRequest, "Id is required")
		log.Error("Error, Id is required")
		return
	}
	if u.Name == "" && u.GroupId <= 0 && u.Date == "" && u.Text == "" && u.Link == "" {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	var date time.Time
	if u.Date != "" {
		date, err = time.Parse("02.01.2006", u.Date)
		if err != nil {
			g.JSON(http.StatusBadRequest, "impt bad date")
			log.Error("Error, impt bad date")
			return
		}
	}

	m := &datastruct.Music{
		Id:      u.Id,
		Name:    u.Name,
		GroupId: u.GroupId,
		Date:    date,
		Text:    u.Text,
		Link:    u.Link,
	}

	err = a.Serv.UpdateMusicById(ctx, m)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadId) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Error, UpdateMusic", err.Error())
		return
	}

	g.JSON(http.StatusOK, m)
}

// DeleteMusic godoc
// @Summary Delete Music
// @Schemes
// @Description do Delete Music
// @Tags DeleteMusic
// @Accept json
// @Produce json
// @Param        m   body      deleteReq  true  "DeleteMusic"
// @Success 200 {object} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 500 {string} ServerError
// @Router /music [delete]
func (a *api) deleteMusicById(g *gin.Context) {
	ctx := g.Request.Context()
	m := deleteReq{}

	err := g.ShouldBindJSON(&m)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body", err.Error())
		return
	}
	if m.Id <= 0 {
		g.JSON(http.StatusBadRequest, "Id is required")
		log.Error("Error, Id is required")
		return
	}

	err = a.Serv.DeleteMusicById(ctx, m.Id)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadId) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Error DeleteMusic", err.Error())
		return
	}

	g.JSON(http.StatusOK, m)
}

// ListMusic godoc
// @Summary ListMusic
// @Schemes
// @Description do List Music
// @Tags ListMusic
// @Accept json
// @Produce json
// @Success 200 {array} datastruct.Music
// @Failure 500 {string} ServerError
// @Router /music/list [get]
func (a *api) getList(g *gin.Context) {
	ctx := g.Request.Context()
	songs, err := a.Serv.GetList(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, songs)
	return
}
