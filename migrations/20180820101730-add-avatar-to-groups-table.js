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
    db.addColumn('chats', 'avatar', {
        type: 'longtext',
        notNull: true,
    }, callback);
};

exports.down = function(db, callback) {
    db.removeColumn('chats', 'avatar', callback);
};
