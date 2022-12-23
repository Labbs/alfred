-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `users` (
  `id` text NOT NULL,
  `username` varchar(255) NOT NULL,
  `email` varchar(255),
  `password` varchar(255) NOT NULL,
  `avatar` text DEFAULT "",
  `dark_mode` text DEFAULT "auto",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `tokens` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `user_id` text NOT NULL,
  `scope` longtext NOT NULL,
  `created_at` integer NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_tokens` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `dashboards` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `default` boolean DEFAULT false,
  `user_id` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_dashboards` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `widgets` (
  `id` text NOT NULL,
  `name` text,
  `x` integer,
  `y` integer,
  `w` integer,
  `h` integer,
  `html` longtext,
  `css` longtext,
  `js` longtext,
  `dashboard_id` text NOT NULL,
  `user_id` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_dashboards_widgets` FOREIGN KEY (`dashboard_id`) REFERENCES `dashboards` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_users_widgets` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `bookmarks` (
  `id` text NOT NULL,
  `name` text,
  `url` text,
  `icon` text,
  `description` text,
  `user_id` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_bookmarks` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `tags` (
  `id` text NOT NULL,
  `name` text NOT NULL,
  `user_id` text NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_users_tags` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `bookmarks_tags` (
  `bookmark_id` text NOT NULL,
  `tag_id` text NOT NULL,
  PRIMARY KEY (`bookmark_id`, `tag_id`),
  CONSTRAINT `fk_bookmarks_bookmarks_tags` FOREIGN KEY (`bookmark_id`) REFERENCES `bookmarks` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_tags_bookmarks_tags` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `bookmarks_tags`;
DROP TABLE IF EXISTS `tags`;
DROP TABLE IF EXISTS `bookmarks`;
DROP TABLE IF EXISTS `widgets`;
DROP TABLE IF EXISTS `dashboards`;
DROP TABLE IF EXISTS `tokens`;
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
