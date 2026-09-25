CREATE DATABASE IF NOT EXISTS pet_social DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE pet_social;

CREATE TABLE IF NOT EXISTS users (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    open_id       VARCHAR(128)    NOT NULL,
    platform      VARCHAR(32)     NOT NULL COMMENT 'alipay|wechat|dev',
    nickname      VARCHAR(64)     NOT NULL DEFAULT '',
    avatar        VARCHAR(512)    NOT NULL DEFAULT '',
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_platform_openid (platform, open_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS pets (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id       BIGINT UNSIGNED NOT NULL,
    name          VARCHAR(64)     NOT NULL,
    level         INT             NOT NULL DEFAULT 1,
    hunger        INT             NOT NULL DEFAULT 50 COMMENT '0-100, lower is hungrier',
    mood          INT             NOT NULL DEFAULT 50 COMMENT '0-100',
    outfit_json   JSON            NOT NULL,
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_user_id (user_id),
    CONSTRAINT fk_pets_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS items (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code          VARCHAR(64)     NOT NULL,
    name          VARCHAR(64)     NOT NULL,
    category      VARCHAR(32)     NOT NULL COMMENT 'outfit|food',
    meta_json     JSON            NULL,
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_items (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id       BIGINT UNSIGNED NOT NULL,
    item_id       BIGINT UNSIGNED NOT NULL,
    quantity      INT             NOT NULL DEFAULT 1,
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_user_item (user_id, item_id),
    CONSTRAINT fk_ui_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_ui_item FOREIGN KEY (item_id) REFERENCES items(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS friendships (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id       BIGINT UNSIGNED NOT NULL,
    friend_id     BIGINT UNSIGNED NOT NULL,
    status        VARCHAR(16)     NOT NULL DEFAULT 'pending' COMMENT 'pending|accepted',
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    UNIQUE KEY uk_pair (user_id, friend_id),
    KEY idx_friend (friend_id),
    CONSTRAINT fk_fs_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_fs_friend FOREIGN KEY (friend_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS chat_messages (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    channel       VARCHAR(16)     NOT NULL COMMENT 'world|friend',
    from_id       BIGINT UNSIGNED NOT NULL,
    to_id         BIGINT UNSIGNED NULL,
    content       VARCHAR(512)    NOT NULL,
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    KEY idx_world (channel, id),
    KEY idx_dm (channel, from_id, to_id, id),
    CONSTRAINT fk_cm_from FOREIGN KEY (from_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS visit_logs (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    visitor_id    BIGINT UNSIGNED NOT NULL,
    host_id       BIGINT UNSIGNED NOT NULL,
    action        VARCHAR(32)     NOT NULL COMMENT 'view|pet|gift',
    created_at    DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    KEY idx_host (host_id, created_at),
    CONSTRAINT fk_vl_visitor FOREIGN KEY (visitor_id) REFERENCES users(id),
    CONSTRAINT fk_vl_host FOREIGN KEY (host_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO items (code, name, category, meta_json) VALUES
    ('hat_basic', '基础帽子', 'outfit', '{"slot":"hat"}'),
    ('hat_star', '星星帽', 'outfit', '{"slot":"hat"}'),
    ('clothes_basic', '基础衣服', 'outfit', '{"slot":"clothes"}'),
    ('clothes_cape', '小披风', 'outfit', '{"slot":"clothes"}'),
    ('food_kibble', '普通狗粮', 'food', '{"hunger":20,"mood":5}'),
    ('food_treat', '美味零食', 'food', '{"hunger":35,"mood":15}')
ON DUPLICATE KEY UPDATE name=VALUES(name);
