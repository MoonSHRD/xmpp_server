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

exports.up = function(db) {
  db.dropTable("offline_messages");
  db.dropTable("chats_msgs");

  db.changeColumn("chats_users", "admin", {type: "varchar", length:15});
  db.renameColumn("chats_users", "admin", "role");
  return null;
};

exports.down = function(db) {
    db.createTable("offline_messages", {
        username:       {type:"varchar", notNull:true},
        data:           {type:"text", notNull:true},
        created_at:     {type:"date", notNull:true}
    });
    db.createTable("chats_msgs", {
        id:             {type: "integer", length:11, primaryKey:true, autoIncrement:true},
        chat_id:        {type:"integer", length:11, notNull:true},
        username:       {type:"varchar", length:256, notNull: true},
        admin:          {type:"integer", length:1, notNull:true, default:0},
        created_at:     {type:"datetime", notNull:true},
        updated_at:     {type:"datetime", notNull:true}
    });

    db.changeColumn("chats_users", "role", {name: "admin", type: "integer", length:1});
    db.renameColumn("chats_users", "role", "admin");
    return null;
};

exports._meta = {
  "version": 1
};
