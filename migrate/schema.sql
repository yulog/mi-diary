CREATE TABLE "notes" ("id" VARCHAR NOT NULL, "user_id" VARCHAR NOT NULL, "reaction_name" VARCHAR, PRIMARY KEY ("id", "user_id"), FOREIGN KEY ("reaction_name") REFERENCES "reactions" ("name") ON UPDATE NO ACTION ON DELETE NO ACTION, FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
CREATE TABLE "users" ("id" VARCHAR NOT NULL, "name" VARCHAR, "count" INTEGER, PRIMARY KEY ("id"));
CREATE TABLE "reactions" ("name" VARCHAR NOT NULL, "image" VARCHAR, "count" INTEGER, PRIMARY KEY ("name"));
CREATE TABLE "hash_tags" ("id" INTEGER NOT NULL, "text" VARCHAR, PRIMARY KEY ("id"), UNIQUE ("text"));
CREATE TABLE "note_to_tags" ("note_id" VARCHAR NOT NULL, "hash_tag_id" INTEGER NOT NULL, PRIMARY KEY ("note_id", "hash_tag_id"), FOREIGN KEY ("note_id") REFERENCES "notes" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, FOREIGN KEY ("hash_tag_id") REFERENCES "hash_tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
