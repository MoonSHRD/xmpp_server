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
  db.createTable('users_messages', {
    id:         { type: dataType.INTEGER, primaryKey: true, autoIncrement: true },
    recipient:  { type: dataType.STRING, length: 256, notNull: true},
    sender:     { type: dataType.STRING, length: 256, notNull: true},
    message:    { type: dataType.TEXT, notNull: true},
    created_at: { type: dataType.DATE_TIME, notNull: true},
    updated_at: { type: dataType.DATE_TIME, notNull: true}

  }, callback);
};

exports.down = function(db, callback) {
  db.dropTable('users_messages', callback);
  return null;
};
