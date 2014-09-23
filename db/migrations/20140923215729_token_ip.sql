
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `tokens` ADD ip varbinary(16);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

