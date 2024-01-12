alter table `user`
    add password_update_time datetime null after password;

update `user` set password_update_time = last_login_time;