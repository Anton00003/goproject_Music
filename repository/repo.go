package repository

import (
	"context"
	"database/sql"
	"fmt"
	"goproject_Music/datastruct"

	//	"reflect"
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
	log.Debug("Repository: ping was successful, New Repository created")
	return r, nil
}

func (r *repo) GetMusicById(ctx context.Context, id int) (*datastruct.Music, error) {
	row := r.Database.QueryRowContext(ctx, "SELECT id, name, groupId, date, text, link FROM songs WHERE id = $1", id) //       добавить контекст
	//	rows, err := r.Database.QueryContext(ctx, "SELECT id, name, groupId, date, text, link FROM songs WHERE id = $1", id)
	//	if err != nil {
	//		return nil, errors.WithMessage(err, "Repository: Request error from DB")
	//	}

	log.Debug("Repository: Request was successful")

	m := &datastruct.Music{}
	//	for rows.Next() { // попробовать убрать, все равно будет одна строка
	err := row.Scan(&(m.Id), &(m.Name), &(m.GroupId), &(m.Date), &(m.Text), &(m.Link))
	if err != nil {
		return nil, err
	}
	//	}
	//	if reflect.DeepEqual(m, &datastruct.Music{}) {
	//		return nil, datastruct.ErrBadId
	//	}

	log.Debug("Repository: get Music was successful")

	return m, nil
}

func (r *repo) GetMusicByFilter(ctx context.Context, filter *datastruct.Music, nOnPage, nPage int) ([]datastruct.Music, error) {
	log.Debug("Filter Run Repository")
	filterText := "SELECT id, name, groupId, date, text, link FROM songs"
	filterN := 0
	s := []any{}

	if filter.Id != 0 {
		filterN++
		s = append(s, filter.Id)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE id = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND id = $%v", filterN)
		}
	}
	if filter.Name != "" {
		filterN++
		s = append(s, filter.Name)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE name = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND name = $%v", filterN)
		}
	}
	if filter.GroupId != 0 {
		filterN++
		s = append(s, filter.GroupId)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE groupId = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND groupId = $%v", filterN)
		}
	}
	if !filter.Date.IsZero() {
		filterN++
		s = append(s, filter.Date)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE date = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND date = $%v", filterN)
		}
	}
	if filter.Text != "" {
		filterN++
		s = append(s, filter.Text)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE text = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND text = $%v", filterN)
		}
	}
	if filter.Link != "" {
		filterN++
		s = append(s, filter.Link)
		if filterN == 1 {
			filterText = filterText + fmt.Sprintf(" WHERE link = $%v", filterN)
		} else {
			filterText = filterText + fmt.Sprintf(" AND link = $%v", filterN)
		}
	}

	filterN = filterN + 2
	s = append(s, nOnPage, (nPage-1)*nOnPage)
	filterText = filterText + fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(filterN-1), strconv.Itoa(filterN))
	fmt.Println("filterText=", filterText)
	fmt.Println("s=", s)

	log.Debug("Repository: input fields have been processed")

	songs := []datastruct.Music{}

	rows, err := r.Database.QueryContext(ctx, filterText, s...)
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: Request error from DB")
	}

	for rows.Next() {
		m := datastruct.Music{}
		rows.Scan(&(m.Id), &(m.Name), &(m.GroupId), &(m.Date), &(m.Text), &(m.Link))
		songs = append(songs, m)
	}

	if len(songs) < 1 {
		return nil, datastruct.ErrBadList
	}
	//	if reflect.DeepEqual(songs, []datastruct.Music{}) {
	//		return nil, datastruct.ErrBadFilter
	//	}

	log.Debug("Repository: get Music was successful")
	return songs, nil
}

func (r *repo) AddMusic(ctx context.Context, m *datastruct.Music) error {
	_, err := r.Database.QueryContext(ctx, "INSERT INTO songs (name, groupId, date, text, link) VALUES ($1, $2, $3, $4, $5)", m.Name, m.GroupId, m.Date, m.Text, m.Link)
	if err != nil {
		log.Debug("Repository: Request error from DB when adding record: ", err)
		return errors.WithMessage(err, "Repository: Request error from DB when adding record")
	}
	log.Debug("Repository: add Music was successful")
	return nil
}

func (r *repo) UpdateMusicById(ctx context.Context, m *datastruct.Music) error {
	log.Debug("Update Run Repository")
	updateText := "UPDATE songs SET "
	updateN := 0
	s := []any{}

	if m.Name != "" {
		updateN++
		s = append(s, m.Name)
		if updateN == 1 {
			updateText = updateText + fmt.Sprintf("name = $%v", updateN)
		} else {
			updateText = updateText + fmt.Sprintf(", name = $%v", updateN)
		}
	}

	if m.GroupId != 0 {
		_, err := r.GetGroupName(ctx, m.GroupId)
		if err != nil {
			return datastruct.ErrBadGroupId
		}
		updateN++
		s = append(s, m.GroupId)
		if updateN == 1 {
			updateText = updateText + fmt.Sprintf("groupId = $%v", updateN)
		} else {
			updateText = updateText + fmt.Sprintf(", groupId = $%v", updateN)
		}
	}
	if !m.Date.IsZero() {
		updateN++
		s = append(s, m.Date)
		if updateN == 1 {
			updateText = updateText + fmt.Sprintf("date = $%v", updateN)
		} else {
			updateText = updateText + fmt.Sprintf(", date = $%v", updateN)
		}
	}
	if m.Text != "" {
		updateN++
		s = append(s, m.Text)
		if updateN == 1 {
			updateText = updateText + fmt.Sprintf("text = $%v", updateN)
		} else {
			updateText = updateText + fmt.Sprintf(", text = $%v", updateN)
		}
	}
	if m.Link != "" {
		updateN++
		s = append(s, m.Link)
		if updateN == 1 {
			updateText = updateText + fmt.Sprintf("link = $%v", updateN)
		} else {
			updateText = updateText + fmt.Sprintf(", link = $%v", updateN)
		}
	}

	updateN++
	s = append(s, m.Id)
	updateText = updateText + fmt.Sprintf(" WHERE id = $%v", updateN)

	fmt.Println("updateText=", updateText)
	fmt.Println("s=", s)

	log.Debug("Repository: input fields have been processed")

	result, err := r.Database.ExecContext(ctx, updateText, s...)
	if err != nil {
		return errors.WithMessage(err, "Repository: Request error from DB when updating record")
	}

	affectedRows, _ := result.RowsAffected()
	fmt.Println("affectedRows = ", affectedRows)
	if affectedRows < 1 {
		return datastruct.ErrBadId
	}

	log.Debug("Repository: update Music was successful")
	return nil
}

func (r *repo) DeleteMusicById(ctx context.Context, id int) error {
	result, err := r.Database.ExecContext(ctx, "DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		log.Debug("Repository: Request error from DB when deleting record: ", err)
		return errors.WithMessage(err, "Repository: Request error from DB when deleting record")
	}

	affectedRows, _ := result.RowsAffected()
	fmt.Println("affectedRows = ", affectedRows)
	if affectedRows < 1 {
		return datastruct.ErrBadId
	}

	log.Debug("Repository: dalete Music was successful")
	return nil
}

func (r *repo) GetGroupId(ctx context.Context, group string) (int, error) {
	rows, err := r.Database.QueryContext(ctx, "SELECT id FROM groups WHERE name = $1", group)
	if err != nil {
		return 0, errors.WithMessage(err, "Repository: Request error from DB")
	}

	log.Debug("Repository: Request was successful")

	var groupId int

	for rows.Next() {
		rows.Scan(&groupId)
		log.Debug("Repository: Row scan, groupId = ", groupId)
	}
	if groupId == 0 {
		return 0, datastruct.ErrBadGroup
	}

	log.Debug("Repository: get GroupId was successful")

	return groupId, nil
}

func (r *repo) AddGroupId(ctx context.Context, group string) error {
	_, err := r.Database.QueryContext(ctx, "INSERT INTO groups (name) VALUES ($1)", group)
	if err != nil {
		log.Debug("Repository: Request error from DB when adding record: ", err)
		return errors.WithMessage(err, "Repository: Request error from DB when adding record")
	}
	log.Debug("Repository: add Group was successful")
	return nil
}

func (r *repo) GetGroupName(ctx context.Context, id int) (string, error) {
	rows, err := r.Database.QueryContext(ctx, "SELECT name FROM groups WHERE id = $1", id)
	if err != nil {
		return "", errors.WithMessage(err, "Repository: Request error from DB")
	}

	log.Debug("Repository: Request was successful")

	var group string

	for rows.Next() {
		rows.Scan(&group)
		log.Debug("Repository: Row scan, group = ", group)
	}
	if group == "" {
		return "", datastruct.ErrBadGroupId
	}

	log.Debug("Repository: get GroupId was successful")

	return group, nil
}

func (r *repo) GetList(ctx context.Context) ([]datastruct.MusicListItem, error) {
	songs := []datastruct.MusicListItem{}

	rows, err := r.Database.QueryContext(ctx, "SELECT songs.id, songs.name, groupId, date, text, link, groups.name FROM songs JOIN groups ON songs.groupId = groups.id")
	if err != nil {
		return nil, errors.WithMessage(err, "Repository: Request error from DB")
	}

	for rows.Next() {
		m := datastruct.MusicListItem{}
		rows.Scan(&(m.Id), &(m.Name), &(m.GroupId), &(m.Date), &(m.Text), &(m.Link), &(m.Group))
		songs = append(songs, m)
	}

	if len(songs) < 1 {
		return nil, datastruct.ErrBadList
	}

	log.Debug("Repository: get Music was successful")
	return songs, nil
}
