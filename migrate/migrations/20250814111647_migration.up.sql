-- create index "users_user_id" to table: "users"
CREATE UNIQUE INDEX `users_user_id` ON `users` (`user_id`);
-- create index "notes_note_id" to table: "notes"
CREATE UNIQUE INDEX `notes_note_id` ON `notes` (`note_id`);
-- create index "files_file_id" to table: "files"
CREATE UNIQUE INDEX `files_file_id` ON `files` (`file_id`);
