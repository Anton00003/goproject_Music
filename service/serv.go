package service

import (
	"context"
	"goproject_Music/datastruct"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type repo interface {
	GetMusicById(ctx context.Context, id int) (*datastruct.Music, error)
	AddMusic(ctx context.Context, m *datastruct.Music) error
	GetMusicByFilter(ctx context.Context, filter *datastruct.Music, nOnPage int, nPage int) ([]datastruct.Music, error)
	UpdateMusicById(ctx context.Context, m *datastruct.Music) error
	DeleteMusicById(ctx context.Context, id int) error
	GetGroupId(ctx context.Context, group string) (int, error)
	AddGroupId(ctx context.Context, group string) error
	GetList(ctx context.Context) ([]datastruct.MusicListItem, error)
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

func (s *serv) GetAllTextMusicById(ctx context.Context, id int) (string, error) {
	m, err := s.Repo.GetMusicById(ctx, id)
	if err != nil {
		return "", err
	}
	return m.Text, nil
}

func (s *serv) GetPaginTextMusicById(ctx context.Context, id, nOnPage, nPage int) ([]string, error) {
	m, err := s.Repo.GetMusicById(ctx, id)
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

func (s *serv) AddMusic(ctx context.Context, song, group string) (*datastruct.Music, error) {
	detail, err := s.GetSongFromClient(song, group)
	if err != nil {
		return nil, err
	}

	var groupId int
	groupId, err = s.GetGroupId(ctx, group)
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

	return m, s.Repo.AddMusic(ctx, m)
}

func (s *serv) GetMusicByFilter(ctx context.Context, filter *datastruct.Music, nOnPage, nPage int) ([]datastruct.Music, error) {
	return s.Repo.GetMusicByFilter(ctx, filter, nOnPage, nPage)
}

func (s *serv) UpdateMusicById(ctx context.Context, m *datastruct.Music) error {
	return s.Repo.UpdateMusicById(ctx, m)
}

func (s *serv) DeleteMusicById(ctx context.Context, id int) error {
	return s.Repo.DeleteMusicById(ctx, id)
}

func (s *serv) GetGroupId(ctx context.Context, group string) (int, error) {
	id, err := s.Repo.GetGroupId(ctx, group)
	if err != nil && !errors.Is(err, datastruct.ErrBadGroup) {
		return 0, err
	}
	if err == nil {
		return id, nil
	}

	err = s.Repo.AddGroupId(ctx, group)
	if err != nil {
		return 0, err
	}

	id, err = s.Repo.GetGroupId(ctx, group)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *serv) GetList(ctx context.Context) ([]datastruct.MusicListItem, error) {
	return s.Repo.GetList(ctx)
}

func (s *serv) GetSongFromClient(name, group string) (*datastruct.SongDetail, error) {
	return s.Client.GetSongFromClient(name, group)
}
