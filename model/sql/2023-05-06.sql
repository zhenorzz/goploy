alter table monitor add success_script text;
alter table monitor add fail_script text ;
alter table monitor add success_server_id int(10) DEFAULT -1;
alter table monitor add fail_server_id int(10) DEFAULT -1 ;
alter table project add tag varchar(60) default '';