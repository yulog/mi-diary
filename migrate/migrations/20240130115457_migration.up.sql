-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_days" table
CREATE TABLE `new_days` (`ymd` varchar NOT NULL, `ym` varchar NULL, `count` integer NULL, PRIMARY KEY (`ymd`), CONSTRAINT `0` FOREIGN KEY (`ym`) REFERENCES `months` (`ym`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- copy rows from old table "days" to new temporary table "new_days"
INSERT INTO `new_days` (`ymd`, `ym`, `count`) SELECT `ymd`, `ym`, `count` FROM `days`;
-- drop "days" table after copying rows
DROP TABLE `days`;
-- rename temporary table "new_days" to "days"
ALTER TABLE `new_days` RENAME TO `days`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
