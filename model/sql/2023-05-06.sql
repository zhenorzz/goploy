alter table monitor add success_script longtext;
alter table monitor add fail_script longtext ;
alter table monitor add success_server_id int(10) default -1;
alter table monitor add fail_server_id int(10) default -1 ;
alter table project add tag varchar(255) default '';