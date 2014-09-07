
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(100) NOT NULL DEFAULT '',
  `password` CHAR(60) NOT NULL DEFAULT '',
  `email` VARCHAR(254) NOT NULL,
  `avatar` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `activated` BOOLEAN NOT NULL DEFAULT 0,
  `activation_token` CHAR(32),
  `password_reset_token` CHAR(32),
  `facebook_id` VARCHAR(128),
  `google_id` VARCHAR(128),
  PRIMARY KEY (`id`),
  UNIQUE KEY (`email`)
) ENGINE = INNODB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `tokens` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `key` CHAR(64) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  UNIQUE KEY (`key`)
) ENGINE = INNODB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `domains` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(63) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  `expired_at` DATETIME NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE,
  UNIQUE KEY (`name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `records` (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
  `name` VARCHAR(63) NOT NULL,
  `type` VARCHAR(8) NOT NULL,
  `value` TEXT NOT NULL,
  `ttl` INT UNSIGNED NOT NULL DEFAULT 0,
  `priority` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `domain_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`domain_id`) REFERENCES domains(`id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

