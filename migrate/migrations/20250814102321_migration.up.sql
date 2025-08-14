-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_users" table
CREATE TABLE `new_users` (`id` integer NOT NULL, `user_id` varchar NOT NULL, `name` varchar NULL, `display_name` varchar NULL, `avatar_url` varchar NULL, `count` integer NULL, PRIMARY KEY (`id`, `user_id`));
-- copy rows from old table "users" to new temporary table "new_users"
INSERT INTO `new_users` (`user_id`, `name`, `display_name`, `avatar_url`, `count`) SELECT `user_id`, `name`, `display_name`, `avatar_url`, `count` FROM `users`;
-- drop "users" table after copying rows
DROP TABLE `users`;
-- rename temporary table "new_users" to "users"
ALTER TABLE `new_users` RENAME TO `users`;
-- create "new_notes" table
CREATE TABLE `new_notes` (`id` integer NOT NULL, `note_id` varchar NOT NULL, `reaction_id` varchar NOT NULL, `user_id` varchar NOT NULL, `reaction_emoji_name` varchar NOT NULL, `text` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`id`, `note_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`reaction_emoji_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "notes" to new temporary table "new_notes"
INSERT INTO `new_notes` (`note_id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at`) SELECT `note_id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at` FROM `notes`;
-- drop "notes" table after copying rows
DROP TABLE `notes`;
-- rename temporary table "new_notes" to "notes"
ALTER TABLE `new_notes` RENAME TO `notes`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
