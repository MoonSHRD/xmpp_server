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

exports.up = async function(db, callback) {
    await db.createTable("files",
        {
            id:             { type: dataType.INTEGER, primaryKey: true, autoIncrement: true},
            message_id:     { type: dataType.INTEGER, notNull:true},
            hash:           { type: "text", length: 40, notNull: true},
            type:           { type: "text", length: 40, notNull: true},
            name:           { type: "text", notNull: true},
        },
        callback);
    await db.addForeignKey("files", "messages", "fk-files-messages", {"message_id": "id"}, {onUpdate: "CASCADE", onDelete: "RESTRICT"}, callback);
    db.addColumn("messages", "files", {type: dataType.INTEGER, length:1, notNull:true}, callback);
    return null;
};

exports.down = async function(db, callback) {
    // db.removeColumn("messages", "file", callback);
    await db.dropTable("files", callback);
    await db.removeForeignKey("files", "fk-files-messages", callback);
    db.removeColumn("messages", "files", callback);

    return null;
};

exports._meta = {
    "version": 1
};
