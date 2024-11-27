package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goproject_Music/datastruct"
	"reflect"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type repo struct {
	Database *sql.DB
}

func NewRepo(dsn string) (*repo, error) {
	r := &repo{}

	var err error
	r.Database, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: opening Error DB")
	}
	log.Debug("Repository: success opening DB")

	if err := r.Database.Ping(); err != nil {
		return nil, errors.WithMessage(err, "Repository: pinging DB")
	}
	log.Debug("Repository: ping was saccessful, New Repository created")
	return r, nil
}

func (r *repo) GetMusicByNameGroup(name, group string) (*datastruct.Music, error) {
	rows, err := r.Database.Query("SELECT * FROM music WHERE musicName = $1 AND musicGroup = $2", name, group)
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: Request error from DB")
	}

	log.Debug("Repository: Request was successful")

	m := &datastruct.Music{}
	var b []byte

	for rows.Next() {
		rows.Scan(&(m.MusicName), &(m.MusicGroup), &(m.MusicDate), &(m.MusicText), &(m.MusicLink), &b)
		err := json.Unmarshal(b, &(m.MusicTextCouplet))
		if err != nil {
			log.Debug("Repository: Error unmarshal JSON for couplets: ", err)
			return nil, errors.WithMessage(err, "Repository: Error unmarshal JSON for couplets")
		}
	}
	if reflect.DeepEqual(m, &datastruct.Music{}) {
		return nil, datastruct.ErrBadNameGroup
	}

	log.Debug("Repository: get Music was saccessful")

	return m, nil
}

func (r *repo) GetMusicByFilter(filter *datastruct.Music, nOnPage, nPage int) ([]datastruct.Music, error) {
	log.Debug("Filter Run Repository")
	filterText := "SELECT *  FROM music"
	filterN := 0
	s := []any{}

	if filter.MusicName != "" {
		filterN++
		s = append(s, filter.MusicName)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE musicName = $%s", strconv.Itoa(filterN))
		} else {
			filterText = filterText + fmt.Sprintf(" AND musicName = $%s", strconv.Itoa(filterN))
		}
	}
	if filter.MusicGroup != "" {
		filterN++
		s = append(s, filter.MusicGroup)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE musicGroup = $%s", strconv.Itoa(filterN))
		} else {
			filterText = filterText + fmt.Sprintf(" AND musicGroup = $%s", strconv.Itoa(filterN))
		}
	}
	if filter.MusicDate != "" {
		filterN++
		s = append(s, filter.MusicDate)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE musicDate = $%s", strconv.Itoa(filterN))
		} else {
			filterText = filterText + fmt.Sprintf(" AND musicDate = $%s", strconv.Itoa(filterN))
		}
	}
	if filter.MusicText != "" {
		filterN++
		s = append(s, filter.MusicText)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE musicText = $%s", strconv.Itoa(filterN))
		} else {
			filterText = filterText + fmt.Sprintf(" AND musicText = $%s", strconv.Itoa(filterN))
		}
	}
	if filter.MusicLink != "" {
		filterN++
		s = append(s, filter.MusicLink)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE musicLink = $%s", strconv.Itoa(filterN))
		} else {
			filterText = filterText + fmt.Sprintf(" AND musicLink = $%s", strconv.Itoa(filterN))
		}
	}

	filterN = filterN + 2
	s = append(s, nOnPage, (nPage-1)*nOnPage)
	filterText = filterText + fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(filterN-1), strconv.Itoa(filterN))
	fmt.Println("filterText=", filterText)
	fmt.Println("s=", s)

	log.Debug("Repository: input fields have been processed")

	musics := []datastruct.Music{}

	rows, err := r.Database.Query(filterText, s...)
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: Request error from DB")
	}

	var b []byte
	for rows.Next() {
		m := datastruct.Music{}
		rows.Scan(&(m.MusicName), &(m.MusicGroup), &(m.MusicDate), &(m.MusicText), &(m.MusicLink), &b)
		err := json.Unmarshal(b, &(m.MusicTextCouplet))
		if err != nil {
			return nil, errors.WithMessage(err, "Repository: Error unmarshal JSON for couplets")
		}
		musics = append(musics, m)
	}

	if reflect.DeepEqual(musics, []datastruct.Music{}) {
		return nil, datastruct.ErrBadFilter
	}

	log.Debug("Repository: get Music was saccessful")
	return musics, nil
}

func (r *repo) AddMusic(m *datastruct.Music) error {
	b, err := json.Marshal(m.MusicTextCouplet)
	if err != nil {
		return errors.WithMessage(err, "Repository: Error marshal JSON for couplets")
	}

	_, err = r.Database.Query("INSERT INTO music (musicName, musicGroup, musicDate, musicText, musicLink, musicTextCouplet) VALUES ($1, $2, $3, $4, $5, $6)", m.MusicName, m.MusicGroup, m.MusicDate, m.MusicText, m.MusicLink, b)
	if err != nil {
		log.Debug("Repository: Request error from DB when adding record: ", err)
		return errors.WithMessage(err, "Repository: Request error from DB when adding record")
	}
	log.Debug("Repository: add Music was saccessful")
	return nil
}

func (r *repo) UpdateMusicFieldValueByNameGroup(name, group, field, value string) error {
	_, err := r.Database.Query(fmt.Sprintf("UPDATE music SET %s = $1 WHERE musicName = $2 AND musicGroup = $3", field), value, name, group)
	if err != nil {
		return errors.WithMessage(err, "Repository: Request error from DB when updating record")
	}
	log.Debug("Repository: update Music was saccessful")
	return nil
}

func (r *repo) UpdateMusicTextCoupletByNameGroup(name, group, field string, couplets []string) error {
	b, err := json.Marshal(couplets)
	if err != nil {
		return errors.WithMessage(err, "Repository: Error marshal JSON for couplets")
	}
	_, err = r.Database.Query(fmt.Sprintf("UPDATE music SET %s = $1 WHERE musicName = $2 AND musicGroup = $3", field), b, name, group)
	if err != nil {
		return errors.WithMessage(err, "Repository: Request error from DB when updating record")
	}
	log.Debug("Repository: update Music was saccessful")
	return nil
}

func (r *repo) DeleteMusicByNameGroup(name, group string) error {
	_, err := r.Database.Query("DELETE FROM music WHERE musicName = $1 AND musicGroup = $2", name, group)
	if err != nil {
		log.Debug("Repository: Request error from DB when deleting record: ", err)
		return errors.WithMessage(err, "Repository: Request error from DB when deleting record")
	}
	log.Debug("Repository: dalete Music was saccessful")
	return nil
}
