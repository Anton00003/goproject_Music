package service

import "goproject_Music/datastruct"

type repo interface {
	GetMusicByNameGroup(name, group string) (*datastruct.Music, error)
	AddMusic(m *datastruct.Music) error
	GetMusicByFilter(*datastruct.Music, int, int) ([]datastruct.Music, error)
	UpdateMusicFieldValueByNameGroup(name, group, field, valueS string) error
	UpdateMusicTextCoupletByNameGroup(name, group, field string, couplets []string) error
	DeleteMusicByNameGroup(name, group string) error
}

type serv struct {
	Repo repo
}

func NewServ(r repo) *serv {
	return &serv{Repo: r}
}

func (s *serv) GetAllTextMusicByNameGroup(name, group string) (string, error) {
	m, err := s.Repo.GetMusicByNameGroup(name, group)
	if err != nil {
		return "", err
	}
	return m.MusicText, nil
}

func (s *serv) GetPaginTextMusicByNameGroup(name, group string, nOnPage, nPage int) ([]string, error) {
	m, err := s.Repo.GetMusicByNameGroup(name, group)
	if err != nil {
		return []string{}, err
	}

	couplets := m.MusicTextCouplet
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

func (s *serv) AddMusic(m *datastruct.Music) error {
	return s.Repo.AddMusic(m)
}

func (s *serv) GetMusicByFilter(filter *datastruct.Music, nOnPage, nPage int) ([]datastruct.Music, error) {
	return s.Repo.GetMusicByFilter(filter, nOnPage, nPage)
}

func (s *serv) GetCouplet(name string, group string, nCouplet int) (string, error) {
	m, err := s.Repo.GetMusicByNameGroup(name, group)
	if err != nil {
		return "", err
	}
	if nCouplet > len(m.MusicTextCouplet) || nCouplet < 1 {
		return "", datastruct.ErrBadNumCouplet
	}
	textCouplet := m.MusicTextCouplet[nCouplet-1]

	return textCouplet, nil
}

func (s *serv) UpdateMusicFieldValueByNameGroup(name, group, field string, value any) error {
	_, ok := datastruct.TableField[field]
	if !ok {
		return datastruct.ErrBadField
	}
	if field == "musicTextCouplet" {
		var couplets []string
		couplets, ok := value.([]string)
		if !ok {
			return datastruct.ErrBadTypeVal
		}
		return s.Repo.UpdateMusicTextCoupletByNameGroup(name, group, field, couplets)
	}
	valueS, ok := value.(string)
	if !ok {
		return datastruct.ErrBadTypeVal
	}
	return s.Repo.UpdateMusicFieldValueByNameGroup(name, group, field, valueS)
}

func (s *serv) DeleteMusicByNameGroup(name, group string) error {
	return s.Repo.DeleteMusicByNameGroup(name, group)
}

func (s *serv) GetInfoMusicByNameGroup(name, group string) (*datastruct.SongDetail, error) {
	m, err := s.Repo.GetMusicByNameGroup(name, group)
	if err != nil {
		return nil, err
	}
	detail := &datastruct.SongDetail{
		ReleaseDate: m.MusicDate,
		Text:        m.MusicText,
		Link:        m.MusicLink,
	}
	return detail, nil
}
