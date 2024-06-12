CREATE TABLE IF NOT EXISTS `notification_template` (
    `id` int unsigned auto_increment,
    `type`  tinyint unsigned default 0 not null,
    `use_by` varchar(255) default '' not null,
    `title`  varchar(255) default ''  not null,
    `template` text not null,
    `insert_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci;

INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (86, 12, 'ShowNotificationPage',  0, '');
INSERT IGNORE INTO `permission`(`id`, `pid`, `name`, `sort`, `description`) VALUES (87, 12, 'EditNotification',  0, '');

UPDATE `permission` set `name` = 'Setting' WHERE `id` = 7;
UPDATE `permission` set `pid` = 7 WHERE `pid` = 12;

DELETE FROM `permission` WHERE `id` = 12;

INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (1, 1, 'deploy', '{{ .Project.Name }}', 'Deploy: <font color="warning">{{ .Project.Name }}</font>
Publisher: <font color="comment">{{ .Project.PublisherName }}</font>
Author: <font color="comment">{{ .CommitInfo.Author }}</font>
{{ if ne .CommitInfo.Tag }}Tag: <font color="comment">{{ .CommitInfo.Tag }}</font>{{ end }}
Branch: <font color="comment">{{ .CommitInfo.Branch }}</font>
CommitSHA: <font color="comment">{{ .CommitInfo.Commit }}</font>
CommitMessage: <font color="comment">{{ .CommitInfo.Message }}</font>
ServerList:<font color="comment">
{{- range .ProjectServers}}
  {{- if ne .Server.Name .Server.IP}}
    {{- .Server.Name}}({{.Server.IP}})
  {{- else}}
    {{- .Server.IP}}
  {{- end}}
{{- end}}
</font>
{{- if eq .DeployState 2 }}
State: <font color="green">success</font>
{{- else }}
State: <font color="red">fail</font>
{{- end }}
{{- if ne .DeployDetail ""}}
Detail: <font color="comment">{{.DeployDetail}}</font>
{{- end }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (2, 2, 'deploy', '{{ .Project.Name }}', '#### Deploy：{{ .Project.Name }}
#### Publisher：{{ .Project.PublisherName }}
#### Author：{{ .CommitInfo.Author }}
#### {{ if ne .CommitInfo.Tag }}Tag: {{ .CommitInfo.Tag }}{{ end }}
#### Branch：{{ .CommitInfo.Branch }}
#### CommitSHA：{{ .CommitInfo.Commit }}
#### CommitMessage： {{ .CommitInfo.Message }}
####  ServerList:
{{- range .ProjectServers}}
  {{- if ne .Server.Name .Server.IP}}
    {{- .Server.Name}}({{.Server.IP}})
  {{- else}}
    {{- .Server.IP}}
  {{- end}}
{{- end}}
####
{{- if eq .DeployState 2 }}State: <font color="green">success</font>
{{- else }}State: <font color="red">fail</font>
{{- end }}
{{- if ne .DeployDetail ""}}
> Detail: <font color="comment">{{.DeployDetail}}</font>
{{- end }}
');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (3, 3, 'deploy', 'Deploy: {{ .Project.Name }}', 'Publisher: {{ .Project.PublisherName }}
Author: {{ .CommitInfo.Author }}
{{ if ne .CommitInfo.Tag "" }}Tag: {{ .CommitInfo.Tag }}{{ end }}
Branch: {{ .CommitInfo.Branch }}
CommitSHA: {{ .CommitInfo.Commit }}
CommitMessage: {{ .CommitInfo.Message }}
ServerList:
{{- range .ProjectServers}}
  {{- if ne .Server.Name .Server.IP}}
    {{- .Server.Name}}({{.Server.IP}}),
  {{- else}}
    {{- .Server.IP}},
  {{- end}}
{{- end}}
{{- if eq .DeployState 2 }}
State: success
{{- else }}
State: fail
{{- end }}
{{- if ne .DeployDetail ""}}
Detail: {{.DeployDetail }}
{{- end }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (4, 1, 'monitor', '{{ .Monitor.Name }}', 'Monitor: <font color="warning">{{ .Monitor.Name }}</font>
> <font color="warning">can not access</font>
> <font color="comment">{{ .ErrorMsg }}</font>');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (5, 2, 'monitor', '{{ .Monitor.Name }}', '#### Monitor: {{ .Monitor.Name }} can not access
{{ .ErrorMsg }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (6, 3, 'monitor', '{{ .Monitor.Name }}', 'can not access
detail: {{ .ErrorMsg }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (7, 1, 'server_monitor', '{{ .Server.Name }} {{ .ServerMonitor.Item }} Warning', 'Server: {{ .Server.Name }}({{ .Server.Description }})
Item: <font color="warning">{{ .ServerMonitor.Item }} warning</font>
Event: {{ .ServerMonitor.Formula }} value: {{ .CycleValue }}, {{ .ServerMonitor.Operator }} {{ .ServerMonitor.Value }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (8, 2, 'server_monitor', '{{ .Server.Name }} {{ .ServerMonitor.Item }} Warning', 'Server: {{ .Server.Name }}({{ .Server.Description }})
Item: {{ .ServerMonitor.Item }} warning
Event: {{ .ServerMonitor.Formula }} value: {{ .CycleValue }}, {{ .ServerMonitor.Operator }} {{ .ServerMonitor.Value }}');
INSERT INTO goploy.notification_template (id, type, use_by, title, template) VALUES (9, 3, 'server_monitor', '{{ .Server.Name }} {{ .ServerMonitor.Item }} Warning', 'Server: {{ .Server.Name }}({{ .Server.Description }})
Item: {{ .ServerMonitor.Item }} warning
Event: {{ .ServerMonitor.Formula }} value: {{ .CycleValue }}, {{ .ServerMonitor.Operator }} {{ .ServerMonitor.Value }}');
