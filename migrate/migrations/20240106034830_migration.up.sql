-- create "notes" table
CREATE TABLE `notes` (`id` varchar NOT NULL, `user_id` varchar NOT NULL, `reaction_name` varchar NULL, PRIMARY KEY (`id`, `user_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`reaction_name`) REFERENCES `reactions` (`name`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create "users" table
CREATE TABLE `users` (`id` varchar NOT NULL, `name` varchar NULL, `count` integer NULL, PRIMARY KEY (`id`));
-- create "reactions" table
CREATE TABLE `reactions` (`name` varchar NOT NULL, `image` varchar NULL, `count` integer NULL, PRIMARY KEY (`name`));
-- create "hash_tags" table
CREATE TABLE `hash_tags` (`id` integer NOT NULL, `text` varchar NULL, PRIMARY KEY (`id`));
-- create index "hash_tags_text" to table: "hash_tags"
CREATE UNIQUE INDEX `hash_tags_text` ON `hash_tags` (`text`);
-- create "note_to_tags" table
CREATE TABLE `note_to_tags` (`note_id` varchar NOT NULL, `hash_tag_id` integer NOT NULL, PRIMARY KEY (`note_id`, `hash_tag_id`), CONSTRAINT `0` FOREIGN KEY (`hash_tag_id`) REFERENCES `hash_tags` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`note_id`) REFERENCES `notes` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
