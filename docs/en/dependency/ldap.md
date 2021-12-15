# LDAP sign in

The content described in this document belongs to advanced usage functions, involves more technical details, and is suitable for reference by users who are experienced in related functions.

1. Administrator adds members
2. The user uses the username and password in ldap to sign in
3. The system will use the username and password to perform login authentication in the LDAP system. After the login is successful, the user can use the system normally

```
The admin account is not controlled by LDAP
When the ldap.enabled is configured, it means that LDAP is turned on. After that, all account and password systems except admin are only authenticated through LDAP.
```

## Description

|      parameter      |                      Example                       | Required |                                  Description                                  | 
|:-------------------:|:--------------------------------------------------:|:--------:|:-----------------------------------------------------------------------------:|
|       enabled       |                        true                        |    是     |                                   Open ldap                                   |
|         url         |             ldap://10.XXX.XXX.XXX:389/             |    是     |                            The address of the LDAP                            |
|       bindDN        |             cn=admin,dc=example,dc=org             |    否     |              This is needed when LDAP prohibit anonymous access               |
|      password       |                       *****                        |    否     |                           The password of the admin                           |
|       baseDN        |          ou=Members,DC=sensorsdata,DC=cn           |    是     |                                 LDAP Base DN                                  |
|         uid         |                   sAMAccountName                   |    是     |                                   LDAP uid                                    |
|     userFilter      |   (memberOf=cn=xxx,ou=people,dc=example,dc=com)    |    否     | Whether to filter according to the specified filter when searching for users  |

