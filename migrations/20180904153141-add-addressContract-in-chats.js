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
    db.changeColumn('chats', 'id', {type: 'varchar', length: 42, autoIncrement: false, notNull: true}, callback);
    db.renameColumn("messages", "recipient", "chat_id", callback);
};

exports.down = function(db, callback) {
    db.changeColumn('chats', 'id', {type: 'int', length: 11, autoIncrement: true, notNull: false}, callback);
    db.renameColumn("messages", "chat_id", "recipient", callback);
};