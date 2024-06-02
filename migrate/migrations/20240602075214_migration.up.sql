-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_reactions" table
CREATE TABLE `new_reactions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` varchar NOT NULL, `image` varchar NULL, `count` integer NULL);
-- copy rows from old table "reactions" to new temporary table "new_reactions"
INSERT INTO `new_reactions` (`id`, `name`, `image`, `count`) SELECT `id`, `name`, `image`, `count` FROM `reactions`;
-- drop "reactions" table after copying rows
DROP TABLE `reactions`;
-- rename temporary table "new_reactions" to "reactions"
ALTER TABLE `new_reactions` RENAME TO `reactions`;
-- create index "reactions_name" to table: "reactions"
CREATE UNIQUE INDEX `reactions_name` ON `reactions` (`name`);
-- create "new_notes" table
CREATE TABLE `new_notes` (`id` varchar NOT NULL, `reaction_id` varchar NOT NULL, `user_id` varchar NOT NULL, `reaction_emoji_name` varchar NOT NULL, `text` varchar NULL, `created_at` timestamp NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`reaction_emoji_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "notes" to new temporary table "new_notes"
INSERT INTO `new_notes` (`id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at`) SELECT `id`, `reaction_id`, `user_id`, `reaction_emoji_name`, `text`, `created_at` FROM `notes`;
-- drop "notes" table after copying rows
DROP TABLE `notes`;
-- rename temporary table "new_notes" to "notes"
ALTER TABLE `new_notes` RENAME TO `notes`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
