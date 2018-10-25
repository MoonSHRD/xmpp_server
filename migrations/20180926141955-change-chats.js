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
  db.changeColumn("chats", "type", {type:"varchar", length:10, notNull:true}, callback);
  db.changeColumn("chats", "chatname", {type: "varchar", notNull:false}, callback);
  return null;
};

exports.down = function(db, callback) {
    db.changeColumn("chats", "type", {type:"integer", length:1, notNull: false}, callback);
    db.changeColumn("chats", "chatname", {type: "varchar", notNull:false}, callback);
    return null;
};

exports._meta = {
  "version": 1
};
