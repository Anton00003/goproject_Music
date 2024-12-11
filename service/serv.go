package service

import (
	"goproject_Music/datastruct"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type repo interface {
	GetMusicById(id int) (*datastruct.Music, error)
	AddMusic(m *datastruct.Music) error
	GetMusicByFilter(*datastruct.Music, int, int) ([]datastruct.Music, error)
	UpdateMusicById(*datastruct.Music) error
	DeleteMusicById(id int) error
	GetGroupId(group string) (int, error)
	AddGroupId(group string) error
	GetList() ([]datastruct.MusicListItem, error)
}

type client interface {
	GetSongFromClient(name, group string) (*datastruct.SongDetail, error)
}

type serv struct {
	Repo   repo
	Client client
}

func NewServ(r repo, c client) *serv {
	return &serv{Repo: r, Client: c}
}

func (s *serv) GetAllTextMusicById(id int) (string, error) {
	m, err := s.Repo.GetMusicById(id)
	if err != nil {
		return "", err
	}
	return m.Text, nil
}

func (s *serv) GetPaginTextMusicById(id, nOnPage, nPage int) ([]string, error) {
	m, err := s.Repo.GetMusicById(id)
	if err != nil {
		return []string{}, err
	}

	couplets := strings.Split(m.Text, "\n\n")
	var coupOnPage []string

	if len(couplets) <= (nPage-1)*nOnPage {
		return []string{}, datastruct.ErrBadNumPage
	}
	if len(couplets) <= (nPage)*nOnPage {
		coupOnPage = couplets[(nPage-1)*nOnPage:]
	} else {
		coupOnPage = couplets[(nPage-1)*nOnPage : nPage*nOnPage]
	}

	return coupOnPage, nil
}

func (s *serv) AddMusic(song, group string) (*datastruct.Music, error) {
	detail, err := s.GetSongFromClient(song, group)
	if err != nil {
		return nil, err
	}

	var groupId int
	groupId, err = s.GetGroupId(group)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("02.01.2006", detail.ReleaseDate) // вынести в контанты
	if err != nil {
		return nil, err
	}

	m := &datastruct.Music{
		Name:    song,
		GroupId: groupId,
		Date:    date,
		Text:    detail.Text,
		Link:    detail.Link,
	}

	return m, s.Repo.AddMusic(m)
}

func (s *serv) GetMusicByFilter(filter *datastruct.Music, nOnPage, nPage int) ([]datastruct.Music, error) {
	return s.Repo.GetMusicByFilter(filter, nOnPage, nPage)
}

func (s *serv) UpdateMusicById(m *datastruct.Music) error {
	return s.Repo.UpdateMusicById(m)
}

func (s *serv) DeleteMusicById(id int) error {
	return s.Repo.DeleteMusicById(id)
}

func (s *serv) GetGroupId(group string) (int, error) {
	id, err := s.Repo.GetGroupId(group)
	if err != nil && !errors.Is(err, datastruct.ErrBadGroup) {
		return 0, err
	}
	if err == nil {
		return id, nil
	}

	err = s.Repo.AddGroupId(group)
	if err != nil {
		return 0, err
	}

	id, err = s.Repo.GetGroupId(group)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serv) GetList() ([]datastruct.MusicListItem, error) {
	return s.Repo.GetList()
}

func (s *serv) GetSongFromClient(name, group string) (*datastruct.SongDetail, error) {
	return s.Client.GetSongFromClient(name, group)
}
