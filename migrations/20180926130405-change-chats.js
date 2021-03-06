'use strict';

var dbm;
var type;
var seed;

/**
  * We receive the dbmigrate dependency from dbmigrate initially.
  * This enables us to not have to rely on NODE_PATH.
  */
exports.setup = function(options, seedLink) {
  dbm = options.dbmigrate;
  type = dbm.dataType;
  seed = seedLink;
};

exports.up = function(db, callback) {
    db.removeForeignKey("chats_users", "fk-chats_users-chats", callback);
    db.changeColumn("chats", "id", {type: "varchar", length:85}, callback);
    db.changeColumn("chats_users", "chat_id", {type: "varchar", length:85}, callback);
    db.changeColumn("messages", "chat_id", {type: "varchar", length:85}, callback);
    db.addForeignKey("chats_users", "chats", "fk-chats_users-chats", {"chat_id":"id"}, {onUpdate: "CASCADE", onDelete: "RESTRICT"}, callback);
};

exports.down = function(db, callback) {
    db.removeForeignKey("chats_users", "fk-chats_users-chats", callback);
    db.changeColumn("chats", "id", {type: "varchar", length:42}, callback);
    db.changeColumn("chats_users", "chat_id", {type: "varchar", length:42}, callback);
    db.changeColumn("messages", "chat_id", {type: "varchar", length:42}, callback);
    db.addForeignKey("chats_users", "chats", "fk-chats_users-chats", {"chat_id":"id"}, {onUpdate: "CASCADE", onDelete: "RESTRICT"}, callback);
};

exports._meta = {
  "version": 1
};
