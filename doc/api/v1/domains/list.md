# GET /users/:user_id/domains

List domains of a user.

## Request

### Headers

- **Authorization**
  + Optional
  + Format: `token <token>`

### Params

- **user_id**
  + Required

## Response

### 200 (application/json)

Returns a list of domains.

- **id**
  + Type: `int`
- **name**
  + Type: `string`
- **created_at**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **user_id**
  + Type: `int`
- **public**
  + Type: `bool`

### 404

User does not exist.