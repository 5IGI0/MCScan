CREATE TABLE `servers`(
    `id`                INT UNSIGNED    NOT NULL AUTO_INCREMENT,
    `address`           VARCHAR(255)    NOT NULL,
    `normalized_desc`   VARCHAR(255)    NOT NULL,
    `favicon_id`        BIGINT UNSIGNED NOT NULL,
    `first_seen`        INT UNSIGNED    NOT NULL,
    `last_scan`         INT UNSIGNED    NOT NULL,
    `last_success_scan` INT UNSIGNED    NOT NULL,
    `modlist`           TEXT            NOT NULL, -- for simple LIKE queries
    `player_count`      BIGINT UNSIGNED NOT NULL,
    `player_max`        BIGINT UNSIGNED NOT NULL,
    `version_id`        BIGINT UNSIGNED NOT NULL,
    `version_name`      VARCHAR(255)    NOT NULL,
    PRIMARY KEY(`id`)
);

CREATE TABLE `favicons`(
    `id`            BIGINT UNSIGNED NOT NULL, -- first 8 bytes of the favicon's SHA-1 (little endian)
    `raw_favicon`   BLOB            NOT NULL,
    PRIMARY KEY(`id`)
);

CREATE TABLE `mods`(
    `id`        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `server_id` INT UNSIGNED    NOT NULL,
    `modid`     VARCHAR(255)    NOT NULL,
    `version`   VARCHAR(255)    NOT NULL,
    PRIMARY KEY(`id`)
);
CREATE INDEX `idx_server_id` ON `mods`(`server_id`);
CREATE INDEX `idx_modid_ver` ON `mods`(`modid`,`version`);

CREATE TABLE `server_aliases`(
    `id`        INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `server_id` INT UNSIGNED NOT NULL,
    `str_addr`  VARCHAR(255) NOT NULL,
    `int_addr`  INT UNSIGNED, -- for CIDR search
    PRIMARY KEY(`id`)
);
CREATE INDEX `idx_str_addr` ON `server_aliases`(`str_addr`);
CREATE INDEX `idx_int_addr` ON `server_aliases`(`int_addr`);
CREATE INDEX `idx_server_id` ON `server_aliases`(`server_id`);

CREATE TABLE `status_history`(
    `id`        INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `server_id` INT UNSIGNED NOT NULL,
    `at`        INT UNSIGNED NOT NULL,
    `status`    TEXT         NOT NULL,
    `aliases`   TEXT         NOT NULL,
    PRIMARY KEY(`id`)
) ROW_FORMAT=COMPRESSED;
CREATE INDEX `idx_server_id_at` ON `status_history`(`server_id`,`at`);