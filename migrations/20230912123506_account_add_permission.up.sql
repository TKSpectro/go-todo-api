-- modify "accounts" table
ALTER TABLE `accounts` ADD COLUMN `permission` bigint unsigned NULL DEFAULT 0;
