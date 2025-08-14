-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_files" table
CREATE TABLE `new_files` (`id` integer NOT NULL, `file_id` varchar NOT NULL, `name` varchar NULL, `url` varchar NULL, `thumbnail_url` varchar NULL, `type` varchar NULL, `dominant_color` varchar NULL, `group_color` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`id`, `file_id`));
-- copy rows from old table "files" to new temporary table "new_files"
INSERT INTO `new_files` (`file_id`, `name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at`) SELECT `file_id`, `name`, `url`, `thumbnail_url`, `type`, `dominant_color`, `group_color`, `created_at` FROM `files`;
-- drop "files" table after copying rows
DROP TABLE `files`;
-- rename temporary table "new_files" to "files"
ALTER TABLE `new_files` RENAME TO `files`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
