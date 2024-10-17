create sequence if not exists table1_seq;

create table if not exists table1
(
    id          bigint default nextval('table1_seq') primary key,
    n           int8,
    code        varchar(4),
    data        text,
    create_dttm timestamp(6)
);


create table if not exists table2
(
    id   uuid,
    data text
);

COPY table1 (n,code,data,create_dttm) FROM '/data/table1.csv' DELIMITER ',' CSV;
