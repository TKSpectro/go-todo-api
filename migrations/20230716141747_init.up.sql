-- create "accounts" table
CREATE TABLE `accounts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `email` longtext NOT NULL,
  `password` longtext NOT NULL,
  `firstname` longtext NULL,
  `lastname` longtext NULL,
  `token_secret` varchar(8) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_accounts_email` USING HASH (`email`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
-- create "todos" table
CREATE TABLE `todos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `title` longtext NOT NULL,
  `description` longtext NULL,
  `completed` bool NULL DEFAULT 0,
  `account_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_accounts_todos` (`account_id`),
  CONSTRAINT `fk_accounts_todos` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci;
