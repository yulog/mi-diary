-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_users" table
CREATE TABLE `new_users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` varchar NOT NULL, `name` varchar NULL, `display_name` varchar NULL, `avatar_url` varchar NULL, `count` integer NULL);
-- copy rows from old table "users" to new temporary table "new_users"
INSERT INTO `new_users` (`id`, `user_id`, `name`, `display_name`, `avatar_url`, `count`) SELECT `id`, `user_id`, `name`, `display_name`, `avatar_url`, `count` FROM `users`;
-- drop "users" table after copying rows
DROP TABLE `users`;
-- rename temporary table "new_users" to "users"
ALTER TABLE `new_users` RENAME TO `users`;
-- create index "user_id" to table: "users"
CREATE INDEX `user_id` ON `users` (`user_id`);
-- create "new_notes" table
CREATE TABLE `new_notes` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `note_id` varchar NOT NULL, `reaction_id` varchar NOT NULL, `user_id` varchar NOT NULL, `reaction_emoji_name` varchar NOT NULL, `text` varchar NULL, `created_at` timestamp NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`reaction_emoji_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "notes" to new temporary table "new_notes"
INSERT INTO `new_notes` (`id`, `note_id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at`) SELECT `id`, `note_id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at` FROM `notes`;
-- drop "notes" table after copying rows
DROP TABLE `notes`;
-- rename temporary table "new_notes" to "notes"
ALTER TABLE `new_notes` RENAME TO `notes`;
-- create index "note_id" to table: "notes"
CREATE INDEX `note_id` ON `notes` (`note_id`);
-- create "new_files" table
CREATE TABLE `new_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `file_id` varchar NOT NULL, `name` varchar NULL, `url` varchar NULL, `thumbnail_url` varchar NULL, `type` varchar NULL, `dominant_color` varchar NULL, `group_color` varchar NULL, `created_at` timestamp NULL);
-- copy rows from old table "files" to new temporary table "new_files"
INSERT INTO `new_files` (`id`, `file_id`, `name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at`) SELECT `id`, `file_id`, `name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at` FROM `files`;
-- drop "files" table after copying rows
DROP TABLE `files`;
-- rename temporary table "new_files" to "files"
ALTER TABLE `new_files` RENAME TO `files`;
-- create index "file_id" to table: "files"
CREATE INDEX `file_id` ON `files` (`file_id`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
