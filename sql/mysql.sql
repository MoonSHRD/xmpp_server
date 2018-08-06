/*
 * Copyright (c) 2018 Miguel Ángel Ortuño.
 * See the LICENSE file for more information.
 */

CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(256) PRIMARY KEY,
    firstname VARCHAR(256) NOT NULL,
    lastname VARCHAR(256) NOT NULL,
     last_presence TEXT NOT NULL,
     last_presence_at DATETIME NOT NULL,
     updated_at DATETIME NOT NULL,
     created_at DATETIME NOT NULL,
    UNIQUE (username)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roster_notifications (
    contact VARCHAR(256) NOT NULL,
    jid VARCHAR(512) NOT NULL,
    elements TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (contact, jid)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE INDEX i_roster_notifications_jid ON roster_notifications(jid);

CREATE TABLE IF NOT EXISTS roster_items (
    username VARCHAR(256) NOT NULL,
    jid VARCHAR(512) NOT NULL,
    name TEXT NOT NULL,
    subscription TEXT NOT NULL,
    groups TEXT NOT NULL,
    ask BOOL NOT NULL,
    ver INT NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (username, jid)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE INDEX i_roster_items_username ON roster_items(username);
CREATE INDEX i_roster_items_jid ON roster_items(jid);

CREATE TABLE IF NOT EXISTS roster_versions (
    username VARCHAR(256) NOT NULL,
    ver INT NOT NULL DEFAULT 0,
    last_deletion_ver INT NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS blocklist_items (
    username VARCHAR(256) NOT NULL,
    jid VARCHAR(512) NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY(username, jid)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE INDEX i_blocklist_items_username ON blocklist_items(username);

CREATE TABLE IF NOT EXISTS private_storage (
    username VARCHAR(256) NOT NULL,
    namespace VARCHAR(512) NOT NULL,
    data MEDIUMTEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (username, namespace)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE INDEX i_private_storage_username ON private_storage(username);

CREATE TABLE IF NOT EXISTS vcards (
    username VARCHAR(256) PRIMARY KEY,
    vcard MEDIUMTEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS offline_messages (
    username VARCHAR(256) NOT NULL,
    data MEDIUMTEXT NOT NULL,
    created_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE INDEX i_offline_messages_username ON offline_messages(username);

CREATE TABLE IF NOT EXISTS auth_nonce (
  username VARCHAR(256) PRIMARY KEY,
  nonce TEXT NOT NULL,
  created_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chats (
  id int(11) PRIMARY KEY AUTO_INCREMENT,
  chatname VARCHAR(256) NOT NULL,
  creator VARCHAR(256) NOT NULL,
  channel int(1) NOT NULL default 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chats_users (
  id int(11) PRIMARY KEY AUTO_INCREMENT,
  chat_id int(11) NOT NULL,
  username VARCHAR(256) NOT NULL,
  admin int(1) NOT NULL default 0,
  created_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE chats_users ADD UNIQUE chat_user(chat_id, username);

CREATE TABLE IF NOT EXISTS chats_msgs (
  id int(11) PRIMARY KEY AUTO_INCREMENT,
  chat_id int(11) NOT NULL,
  username VARCHAR(256) NOT NULL,
  msg TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
