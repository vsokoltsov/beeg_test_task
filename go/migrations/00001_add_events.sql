-- +goose Up
create table beeg_events (
    id int unsigned not null primary key auto_increment,
    label char(100) not null,
    count int unsigned not null default 1,
    unique key(id, label)
);
-- create unique index id_label on beeg_events (id, label);

-- +goose Down
DROP TABLE beeg_events;
