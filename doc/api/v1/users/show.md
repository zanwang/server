# GET /users/:id

Show a user.

## Request

### Headers

- **Authorization**
  + Optional
  + Format: `token <token>`

### Params

- **id**
  + Required

## Response

### 200 (application/json)

- **id**
  + Type: `int`
- **name**
  + Type: `string`
- **email**
  + Type: `string`
- **avatar**
  + Type: `string`
- **created_at**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **activated**
  + Type: `bool`

### 404

User does not exist
