package api

import (
	"fmt"
	"goproject_Music/datastruct"

	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	docs "goproject_Music/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type updateReq struct {
	Name   string   `json:"name"`
	Group  string   `json:"group"`
	Field  string   `json:"field"`
	Value  string   `json:"value"`
	Values []string `json:"Values"`
}

type deleteReq struct {
	Name  string `json:"name"`
	Group string `json:"group"`
}

type serv interface {
	GetAllTextMusicByNameGroup(name, group string) (string, error)
	GetPaginTextMusicByNameGroup(name string, group string, nOnPage int, nPage int) ([]string, error)
	AddMusic(*datastruct.Music) error
	GetMusicByFilter(*datastruct.Music, int, int) ([]datastruct.Music, error)
	GetCouplet(name string, group string, nCouplet int) (string, error)
	UpdateMusicFieldValueByNameGroup(name, group, field string, value any) error
	DeleteMusicByNameGroup(name, group string) error
	GetInfoMusicByNameGroup(name, group string) (*datastruct.SongDetail, error)
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

	r.GET("/all_text/", a.getAllTextMusicByNameGroup)
	r.GET("/text/", a.getPaginTextMusicByNameGroup)
	r.GET("/couplet/", a.getCouplet)
	r.GET("/info/", a.getInfoMusicByNameGroup)
	r.GET("/music/", a.getMusicByFilter)
	r.POST("/music/", a.addMusic)
	r.PATCH("/music/", a.updateMusicFieldValueByNameGroup)
	r.DELETE("/music/", a.deleteMusicByNameGroup)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(host)
}

// InfoMusic godoc
// @Summary Get InfoMusic
// @Schemes
// @Description do InfoMusic
// @Tags InfoMusic
// @Accept json
// @Produce json
// @Param        name   query      string  false  "Name"
// @Param        group  query      string  false  "Group"
// @Success 200 {object} datastruct.SongDetail
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /info/ [get]
func (a *api) getInfoMusicByNameGroup(g *gin.Context) {
	name := g.Request.URL.Query().Get("name")
	group := g.Request.URL.Query().Get("group")
	fmt.Println("name=", name)
	fmt.Println("group=", group)
	if name == "" || group == "" {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	SongDetail, err := a.Serv.GetInfoMusicByNameGroup(name, group)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadNameGroup) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, *SongDetail)
}

// AllTextMusic godoc
// @Summary Get AllTextMusic
// @Schemes
// @Description do AllTextMusic
// @Tags AllTextMusic
// @Accept json
// @Produce json
// @Param        name   query      string  false  "Name"
// @Param        group   query      string  false  "Group"
// @Success 200 {string} TextSong
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /all_text/ [get]
func (a *api) getAllTextMusicByNameGroup(g *gin.Context) {
	name := g.Request.URL.Query().Get("name")
	group := g.Request.URL.Query().Get("group")

	if name == "" || group == "" {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	TextSong, err := a.Serv.GetAllTextMusicByNameGroup(name, group)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadNameGroup) {
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
// @Param        name     query      string  false  "Name"
// @Param        group    query      string  false  "Group"
// @Param        nOnPage  query      int     false  "nOnPage"
// @Param        nPage    query      int     false  "nPage"
// @Success 200 {string} TextSong
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /text/ [get]
func (a *api) getPaginTextMusicByNameGroup(g *gin.Context) {
	name := g.Request.URL.Query().Get("name")
	group := g.Request.URL.Query().Get("group")
	nOnPageS := g.Request.URL.Query().Get("nOnPage")
	nPageS := g.Request.URL.Query().Get("nPage")

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

	if name == "" || group == "" {
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

	coupOnPage, err := a.Serv.GetPaginTextMusicByNameGroup(name, group, nOnPage, nPage)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadNumPage) {
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
// @Param        name     query      string  false  "Name"
// @Param        group    query      string  false  "Group"
// @Param        date     query      string  false  "Date"
// @Param        text     query      string  false  "Text"
// @Param        link     query      string  false  "Link"
// @Param        nOnPage  query      int     false  "nOnPage"
// @Param        nPage    query      int     false  "nPage"
// @Success 200 {array} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /music/ [get]
func (a *api) getMusicByFilter(g *gin.Context) {
	log.Debug("Api: Filter Run")
	filter := &datastruct.Music{}

	filter.MusicName = g.Request.URL.Query().Get("name")
	filter.MusicGroup = g.Request.URL.Query().Get("group")
	filter.MusicDate = g.Request.URL.Query().Get("date")
	filter.MusicText = g.Request.URL.Query().Get("text")
	filter.MusicLink = g.Request.URL.Query().Get("link")
	nOnPageS := g.Request.URL.Query().Get("nOnPage")
	nPageS := g.Request.URL.Query().Get("nPage")

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
	if filter.MusicName == "" && filter.MusicGroup == "" && filter.MusicDate == "" && filter.MusicText == "" && filter.MusicLink == "" {
		g.JSON(http.StatusBadRequest, "parameters is required")
		log.Error("Error Music parameter")
		return
	}
	if nPage <= 0 {
		g.JSON(http.StatusBadRequest, "page can not be <= 0")
		log.Error("Error Music parameter")
		return
	}

	musics, err := a.Serv.GetMusicByFilter(filter, nOnPage, nPage)
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

	g.JSON(http.StatusOK, musics)
}

// CoupletMusic godoc
// @Summary Get CoupletMusic
// @Schemes
// @Description do CoupletMusic
// @Tags CoupletMusic
// @Accept json
// @Produce json
// @Param        name          query      string  false  "Name"
// @Param        group         query      string  false  "Group"
// @Param        nCouplet      query      string  false  "Couplet"
// @Success 200 {string} TextCouplet
// @Failure 400 {string} BadRequest
// @Failure 404 {string} NotFound
// @Failure 500 {string} ServerError
// @Router /couplet/ [get]
func (a *api) getCouplet(g *gin.Context) {
	name := g.Request.URL.Query().Get("name")
	group := g.Request.URL.Query().Get("group")
	nCoupletStr := g.Request.URL.Query().Get("nCouplet")
	if name == "" || group == "" || nCoupletStr == "" {
		g.JSON(http.StatusBadRequest, "Not enough entered parameters")
		log.Error("Not enough entered parameters")
		return
	}

	nCouplet, err := strconv.Atoi(nCoupletStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, "nCouplet is not numerical")
		log.Error("nCouplet is not numerical: ", err.Error())
		return
	}

	textCouplet, err := a.Serv.GetCouplet(name, group, nCouplet)
	if err != nil {
		if errors.Is(err, datastruct.ErrBadNameGroup) || errors.Is(err, datastruct.ErrBadNumCouplet) {
			g.JSON(http.StatusNotFound, err.Error())
			log.Error("Server Error: ", err.Error())
			return
		}
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Server Error: ", err.Error())
		return
	}

	g.JSON(http.StatusOK, textCouplet)
}

// AddMusic godoc
// @Summary Add Music
// @Schemes
// @Description do Music
// @Tags Music
// @Accept json
// @Produce json
// @Param        m   body      datastruct.Music  true  "Music"
// @Success 200 {object} datastruct.Music
// @Failure 400 {string} BadRequest
// @Failure 500 {string} ServerError
// @Router /music/ [post]
func (a *api) addMusic(g *gin.Context) {
	m := &datastruct.Music{}
	err := g.ShouldBindJSON(m)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body: ", err.Error())
		return
	}

	if m.MusicName == "" {
		g.JSON(http.StatusBadRequest, "name is required")
		log.Error("Error, name is required")
		return
	}
	if m.MusicGroup == "" {
		g.JSON(http.StatusBadRequest, "group is required")
		log.Error("Error, group is required")
		return
	}
	if m.MusicDate == "" {
		g.JSON(http.StatusBadRequest, "date is required")
		log.Error("Error, date is required")
		return
	}
	if m.MusicText == "" {
		g.JSON(http.StatusBadRequest, "text is required")
		log.Error("Error, text is required")
		return
	}
	if m.MusicLink == "" {
		g.JSON(http.StatusBadRequest, "link is required")
		log.Error("Error, link is required")
		return
	}
	if len(m.MusicTextCouplet) == 0 {
		g.JSON(http.StatusBadRequest, "couplet is required")
		log.Error("Error, couplet is required")
		return
	}

	err = a.Serv.AddMusic(m)
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
// @Router /music/ [patch]
func (a *api) updateMusicFieldValueByNameGroup(g *gin.Context) {
	m := updateReq{}

	err := g.ShouldBindJSON(&m)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body", err.Error())
		return
	}
	if m.Name == "" {
		g.JSON(http.StatusBadRequest, "Name is required")
		log.Error("Error, Name is required")
		return
	}
	if m.Group == "" {
		g.JSON(http.StatusBadRequest, "Group is required")
		log.Error("Error, Group is required")
		return
	}
	if m.Field == "" {
		g.JSON(http.StatusBadRequest, "Field is required")
		log.Error("Error, Field is required")
		return
	}
	if m.Value == "" && len(m.Values) == 0 {
		g.JSON(http.StatusBadRequest, "Value/Values is required")
		log.Error("Error, Value is required")
		return
	}
	if m.Value != "" {
		err = a.Serv.UpdateMusicFieldValueByNameGroup(m.Name, m.Group, m.Field, m.Value)
	} else {
		err = a.Serv.UpdateMusicFieldValueByNameGroup(m.Name, m.Group, m.Field, m.Values)
	}
	if err != nil {
		if errors.Is(err, datastruct.ErrBadField) {
			g.JSON(http.StatusBadRequest, err.Error())
			log.Error("Error: ", err.Error())
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
// @Router /music/ [delete]
func (a *api) deleteMusicByNameGroup(g *gin.Context) {
	m := deleteReq{}

	err := g.ShouldBindJSON(&m)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error with read body")
		log.Error("Error with read body", err.Error())
		return
	}
	if m.Name == "" {
		g.JSON(http.StatusBadRequest, "Name is required")
		log.Error("Error, Name is required")
		return
	}
	if m.Group == "" {
		g.JSON(http.StatusBadRequest, "Group is required")
		log.Error("Error, Group is required")
		return
	}

	err = a.Serv.DeleteMusicByNameGroup(m.Name, m.Group)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		log.Error("Error DeleteMusic", err.Error())
		return
	}

	g.JSON(http.StatusOK, m)
}
