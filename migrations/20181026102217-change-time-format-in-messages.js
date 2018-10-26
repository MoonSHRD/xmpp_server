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
  db.changeColumn("messages", "created_at", {type: dataType.BIGINT, length:14, callback});
  db.changeColumn("messages", "updated_at", {type: dataType.BIGINT, length:14, callback});
    return null;
};

exports.down = function(db, callback) {
    db.changeColumn("messages", "created_at", {type: dataType.DATE_TIME}, callback);
    db.changeColumn("messages", "updated_at", {type: dataType.DATE_TIME}, callback);
  return null;
};

exports._meta = {
  "version": 1
};
