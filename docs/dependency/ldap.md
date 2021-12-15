# LDAP 登录

本文档所描述的内容属于高级使用功能，涉及较多技术细节，适用于对相关功能有经验的用户参考。

1. 管理员添加成员
2. 用户在系统中正常使用用户名和密码进行登录
3. 系统会使用用户填写的用户名和密码去后台配置的 LDAP 系统当中进行登录认证，登录成功之后用户就可以正常使用系统

```
admin 帐号不受 LDAP 控制
当配置了 ldap.enabled 配置项，表示开启了 LDAP，此后除了 admin 之外的帐号密码体系，都只通过 LDAP 进行验证。
```

## 配置项说明

|      参数       |                              范例                               |    必填    |                 说明                  | 
|:-------------:|:-------------------------------------------------------------:|:--------:|:-----------------------------------:|
|    enabled    |                             true                              |    是     |               开启ldap                |
|      url      |                  ldap://10.XXX.XXX.XXX:389/                   |    是     |             LDAP 服务的地址              |
|    bindDN     |                  cn=admin,dc=example,dc=org                   |    否     |       当 LDAP 配置了禁止匿名访问的时候需要此项       |
|   password    |                             *****                             |    否     |    当 LDAP 配置了禁止匿名访问的时候需要绑定的管理员密码    |
|    baseDN     |                ou=Members,DC=sensorsdata,DC=cn                |    是     | LDAP Base DN，进行LDAP 用户名检索的 Base Dn  |
|      uid      |                        sAMAccountName                         |    是     |             LDAP 帐号唯一值              |
|  userFilter   |         (memberOf=cn=xxx,ou=people,dc=example,dc=com)         |    否     |  在 LDAP 中查找用户时是否按照指定的 filter 进行筛选   |

