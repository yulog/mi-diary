-- reverse: drop index "file_id" from table: "files"
CREATE INDEX `file_id` ON `files` (`file_id`);
-- reverse: drop index "note_id" from table: "notes"
CREATE INDEX `note_id` ON `notes` (`note_id`);
-- reverse: drop index "user_id" from table: "users"
CREATE INDEX `user_id` ON `users` (`user_id`);
