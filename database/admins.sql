DROP TABLE IF EXISTS `admins`;

CREATE TABLE `admins` (
  `id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_by` varchar(255) NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY (`username`),
  FOREIGN KEY (`created_by`)
    REFERENCES `admins` (`id`)
      ON UPDATE CASCADE
      ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;
