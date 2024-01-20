-- create "months" table
CREATE TABLE `months` (`ym` varchar NOT NULL, `count` integer NULL, PRIMARY KEY (`ym`));
-- create "days" table
CREATE TABLE `days` (`ymd` varchar NOT NULL, `ym` varchar NULL, `count` integer NULL, PRIMARY KEY (`ymd`));
