package main

import (
	"bee01/models"
	"testing"
)

func TestAdd(t *testing.T) {
	models.AddCategory("zhangshan")

}

/*
create table `category`
    -- --------------------------------------------------
    --  Table Structure for `bee01/models.Category`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `category` (
        `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        `title` varchar(255) NOT NULL DEFAULT '' ,
        `created` datetime NOT NULL,
        `views` integer NOT NULL DEFAULT 0 ,
        `topic_time` datetime NOT NULL,
        `topic_count` integer NOT NULL DEFAULT 0 ,
        `topic_last_user_id` integer NOT NULL DEFAULT 0
    );
    CREATE INDEX `category_created` ON `category` (`created`);
    CREATE INDEX `category_views` ON `category` (`views`);
    CREATE INDEX `category_topic_time` ON `category` (`topic_time`);

create table `topic`
    -- --------------------------------------------------
    --  Table Structure for `bee01/models.Topic`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `topic` (
        `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        `uid` integer NOT NULL DEFAULT 0 ,
        `title` varchar(255) NOT NULL DEFAULT '' ,
        `content` varchar(5000) NOT NULL DEFAULT '' ,
        `attachment` varchar(255) NOT NULL DEFAULT '' ,
        `created` datetime NOT NULL,
        `updated` datetime NOT NULL,
        `views` integer NOT NULL DEFAULT 0 ,
        `author` varchar(255) NOT NULL DEFAULT '' ,
        `reply_time` datetime NOT NULL,
        `reply_count` integer NOT NULL DEFAULT 0 ,
        `reply_last_user_id` integer NOT NULL DEFAULT 0
    );
    CREATE INDEX `topic_created` ON `topic` (`created`);
    CREATE INDEX `topic_updated` ON `topic` (`updated`);
    CREATE INDEX `topic_reply_time` ON `topic` (`reply_time`);

*/
