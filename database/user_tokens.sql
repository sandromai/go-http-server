DROP TABLE IF EXISTS `user_tokens`;

CREATE TABLE `user_tokens` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int UNSIGNED NOT NULL,
  `from_login_token` int UNSIGNED NULL,
  `from_user_token` int UNSIGNED NULL,
  `ip_address` varchar(255) NOT NULL,
  `device` varchar(255) NOT NULL,
  `disconnected` boolean NOT NULL DEFAULT false,
  `last_activity` datetime NOT NULL DEFAULT current_timestamp(),
  `expires_at` datetime NOT NULL DEFAULT current_timestamp(),
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY (`from_login_token`),
  FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`)
      ON UPDATE CASCADE
      ON DELETE CASCADE,
  FOREIGN KEY (`from_login_token`)
    REFERENCES `login_tokens` (`id`)
      ON UPDATE CASCADE
      ON DELETE SET NULL,
  FOREIGN KEY (`from_user_token`)
    REFERENCES `user_tokens` (`id`)
      ON UPDATE CASCADE
      ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;
