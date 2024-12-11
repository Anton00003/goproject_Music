-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS groups (id BIGSERIAL PRIMARY KEY, name TEXT);
CREATE UNIQUE INDEX index_group_unique ON groups (name);
CREATE TABLE IF NOT EXISTS songs (id BIGSERIAL PRIMARY KEY, name TEXT, groupId INTEGER, date DATE, text TEXT, link TEXT, FOREIGN KEY (groupId) REFERENCES groups (id));
CREATE UNIQUE INDEX index_songs_unique ON songs (groupId, name);-- можно поставить первым
SELECT 'up SQL query';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
DROP TABLE IF EXISTS groups;
SELECT 'down SQL query';
-- +goose StatementEnd
