-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS music (musicName TEXT, musicGroup TEXT, musicDate TEXT, musicText TEXT, musicLink TEXT, musicTextCouplet JSONB);
CREATE UNIQUE INDEX index_music_unique ON music (musicName, musicGroup);
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE music;
SELECT 'down SQL query';
-- +goose StatementEnd
