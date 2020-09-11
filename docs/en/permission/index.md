# Roles
A user must first allocate a namespace before logging in to the system

All user operations are in the login space, or there are multiple namespaces to switch

## admin 
Have all permissions, do whatever you want

## manager 
Unable to add members, all permissions in the namespace

## group-manager 
Unable to add members and servers, all permissions in the namespace

## member 
Deploy

## Permission code table
| action                              | member | group-manager | manager | admin |
| ------------------------------------| ------ | ------------- | ------- | ----- |
| Deploy                             |   ✓    |       ✓       |    ✓    |   ✓   |
| Monitor                             |        |       ✓       |    ✓    |   ✓   |
| Project                            |        |       ✓       |    ✓    |   ✓   |
| Server                          |        |               |    ✓    |   ✓   |
| Namespace-View、Edit                   |        |               |    ✓    |   ✓   |
| Namespace-Add、Delete                  |        |               |         |   ✓   |
| Member                             |        |               |         |   ✓   |
