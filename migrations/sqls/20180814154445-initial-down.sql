-- DROP INDEX i_roster_notifications_jid ON roster_notifications(jid);
-- DROP INDEX i_roster_items_username ON roster_items(username);
-- DROP INDEX i_roster_items_jid ON roster_items(jid);
-- DROP INDEX i_blocklist_items_username ON blocklist_items(username);
-- DROP INDEX i_private_storage_username ON private_storage(username);
-- DROP INDEX i_offline_messages_username ON offline_messages(username);

DROP TABLE users;
DROP TABLE roster_notifications;
DROP TABLE roster_items;
DROP TABLE roster_versions;
DROP TABLE blocklist_items;
DROP TABLE private_storage;
DROP TABLE offline_messages;
DROP TABLE auth_nonce;
DROP TABLE chats;
DROP TABLE chats_users;
DROP TABLE chats_msgs;
DROP TABLE vcards;