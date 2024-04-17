-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_notes" table
CREATE TABLE `new_notes` (`id` varchar NOT NULL, `user_id` varchar NULL, `reaction_name` varchar NULL, `text` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`reaction_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "notes" to new temporary table "new_notes"
INSERT INTO `new_notes` (`id`, `user_id`, `reaction_name`, `text`, `created_at`) SELECT `id`, `user_id`, `reaction_name`, `text`, `created_at` FROM `notes`;
-- drop "notes" table after copying rows
DROP TABLE `notes`;
-- rename temporary table "new_notes" to "notes"
ALTER TABLE `new_notes` RENAME TO `notes`;
-- create "files" table
CREATE TABLE `files` (`id` varchar NOT NULL, `name` varchar NULL, `url` varchar NULL, `thumbnail_url` varchar NULL, PRIMARY KEY (`id`));
-- create "note_to_files" table
CREATE TABLE `note_to_files` (`note_id` varchar NOT NULL, `file_id` varchar NOT NULL, PRIMARY KEY (`note_id`, `file_id`), CONSTRAINT `0` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`note_id`) REFERENCES `notes` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
