# API v1 Overview

Base URL: `https://maji.moe/api/v1`

## Users

### Object

``` js
{
    "id": 1,
    "name": "John",
    "email": "john@maji.moe",
    "avatar": "//www.gravatar.com/avatar/be3f9c59909ce62242cea65d064bacdc",
    "created_at": "2014-09-10T15:04:05Z",
    "updated_at": "2014-09-10T15:04:05Z",
    "activated": true
}
```

### Methods

- [Create a user](users/create.md)
- [Retrieve a user](users/show.md)
- [Update a user](users/update.md)
- [Delete a user](users/destroy.md)

## Tokens

### Object

``` js
{
    "key": "8dc7deeb918c25fc8160c2467ff1f44db781c465ae5a3a90f613f340a3c6bf77",
    "updated_at": "2014-09-10T15:04:05Z",
    "expired_at": "2014-09-10T15:04:05Z",
    "user_id": 1
}
```

### Methods

- [Create a token](tokens/create.md)
- [Update a token](tokens/update.md)
- [Delete a token](tokens/destroy.md)

## Domains

### Object

``` js
{
    "id": 1,
    "name": "loli",
    "created_at": "2014-09-10T15:04:05Z",
    "updated_at": "2014-09-10T15:04:05Z",
    "expired_at": "2015-09-10T15:04:05Z",
    "user_id": 1
}
```

### Methods

- [Create a domain](domains/create.md)
- [List domains of a user](domains/list.md)
- [Retrieve a domain](domains/show.md)
- [Update a domain](domains/update.md)
- [Delete a domain](domains/destroy.md)

## Records

### Object

``` js
{
    "id": 1,
    "name": "loli",
    "type": "A",
    "value": "127.0.0.1",
    "ttl": 0,
    "priority": 0,
    "created_at": "2014-09-10T15:04:05Z",
    "updated_at": "2014-09-10T15:04:05Z",
    "domain_id": 1
}
```

### Methods

- [Create a record](records/create.md)
- [List records of a domain](records/list.md)
- [Retrieve a record](records/show.md)
- [Update a record](records/update.md)
- [Delete a record](records/destroy.md)