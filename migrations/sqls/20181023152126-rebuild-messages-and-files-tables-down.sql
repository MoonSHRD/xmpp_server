alter table `files` drop FOREIGN KEY `fk-files-messages`;

drop table messages;
drop table files;

create table if not exists messages (
    id int(11) primary key  AUTO_INCREMENT,
    chat_id varchar(85),
    sender varchar(256) not null,
    message TEXT not null,
    created_at DATETIME not null,
    updated_at DATETIME not null,
    delivered int(1) not null,
    files int(1) not null

);

create table if not exists files (
     id int(11) primary key  AUTO_INCREMENT,
     message_id int(11) not null,
     hash TEXT(40) not null,
     type TEXT(40) not null,
     name TEXT not null,
     constraint `fk-files-messages` foreign key (`message_id`) references `messages` (`id`) ON Delete CASCADE on update CASCADE
);
