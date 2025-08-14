-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_users" table
CREATE TABLE `new_users` (`user_id` varchar NOT NULL, `name` varchar NULL, `display_name` varchar NULL, `avatar_url` varchar NULL, `count` integer NULL, PRIMARY KEY (`user_id`));
-- copy rows from old table "users" to new temporary table "new_users"
INSERT INTO `new_users` (`name`, `display_name`, `avatar_url`, `count`) SELECT `name`, `display_name`, `avatar_url`, `count` FROM `users`;
-- drop "users" table after copying rows
DROP TABLE `users`;
-- rename temporary table "new_users" to "users"
ALTER TABLE `new_users` RENAME TO `users`;
-- create "new_note_to_tags" table
CREATE TABLE `new_note_to_tags` (`note_id` varchar NOT NULL, `hash_tag_id` integer NOT NULL, PRIMARY KEY (`note_id`, `hash_tag_id`), CONSTRAINT `0` FOREIGN KEY (`note_id`) REFERENCES `notes` (`note_id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`hash_tag_id`) REFERENCES `hash_tags` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "note_to_tags" to new temporary table "new_note_to_tags"
INSERT INTO `new_note_to_tags` (`note_id`, `hash_tag_id`) SELECT `note_id`, `hash_tag_id` FROM `note_to_tags`;
-- drop "note_to_tags" table after copying rows
DROP TABLE `note_to_tags`;
-- rename temporary table "new_note_to_tags" to "note_to_tags"
ALTER TABLE `new_note_to_tags` RENAME TO `note_to_tags`;
-- create "new_files" table
CREATE TABLE `new_files` (`file_id` varchar NOT NULL, `name` varchar NULL, `url` varchar NULL, `thumbnail_url` varchar NULL, `type` varchar NULL, `dominant_color` varchar NULL, `group_color` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`file_id`));
-- copy rows from old table "files" to new temporary table "new_files"
INSERT INTO `new_files` (`name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at`) SELECT `name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at` FROM `files`;
-- drop "files" table after copying rows
DROP TABLE `files`;
-- rename temporary table "new_files" to "files"
ALTER TABLE `new_files` RENAME TO `files`;
-- create "new_note_to_files" table
CREATE TABLE `new_note_to_files` (`note_id` varchar NOT NULL, `file_id` varchar NOT NULL, PRIMARY KEY (`note_id`, `file_id`), CONSTRAINT `0` FOREIGN KEY (`note_id`) REFERENCES `notes` (`note_id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "note_to_files" to new temporary table "new_note_to_files"
INSERT INTO `new_note_to_files` (`note_id`, `file_id`) SELECT `note_id`, `file_id` FROM `note_to_files`;
-- drop "note_to_files" table after copying rows
DROP TABLE `note_to_files`;
-- rename temporary table "new_note_to_files" to "note_to_files"
ALTER TABLE `new_note_to_files` RENAME TO `note_to_files`;
-- create "new_notes" table
CREATE TABLE `new_notes` (`note_id` varchar NOT NULL, `reaction_id` varchar NOT NULL, `user_id` varchar NOT NULL, `reaction_emoji_name` varchar NOT NULL, `text` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`note_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`reaction_emoji_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "notes" to new temporary table "new_notes"
INSERT INTO `new_notes` (`reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at`) SELECT `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at` FROM `notes`;
-- drop "notes" table after copying rows
DROP TABLE `notes`;
-- rename temporary table "new_notes" to "notes"
ALTER TABLE `new_notes` RENAME TO `notes`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
