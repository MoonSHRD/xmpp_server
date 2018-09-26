'use strict';

var dbm;
var type;
var seed;
var dataType = require('db-migrate-shared').dataType;

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
    db.addColumn("chats", "bio", {type:"text"}, callback);
    db.addColumn("chats", "type", {type:"integer", length:1, notNull:true}, callback);
    db.addColumn("chats", "role", {type:"integer", length:1}, callback);
    db.addColumn("messages", "delivered", {type: "integer", length:1, notNull:true}, callback);

    db.changeColumn("chats_users", "chat_id", {type:"varchar", length:42, notNull:true}, callback);
    db.changeColumn("messages", "chat_id", {type:"varchar", length:42, notNull:true}, callback);

    db.removeColumn("chats", "channel", callback);
    return null;
};

exports.down = function(db, callback) {
    db.removeColumn("chats", "bio", callback);
    db.removeColumn("chats", "type", callback);
    db.removeColumn("chats", "role", callback);
    db.removeColumn("messages", "delivered", callback);

    db.changeColumn("chats_users", "chat_id", {type:"int", length:11}, callback);
    db.changeColumn("messages", "chat_id", {type:"string", length:256, primaryKey:false}, callback);

    db.addColumn("chats", "channel", {type:dataType.INTEGER, length:1, default:'0'}, callback);
    return null;

};

exports._meta = {
  "version": 1
};
