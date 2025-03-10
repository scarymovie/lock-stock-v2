begin;
alter table players add column uid varchar(255) not null default md5(random()::text || now()::text);