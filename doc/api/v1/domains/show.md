# GET /domains/:id

Show a domain.

## Request

### Headers

- **Authorization**
  + Required
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
- **created_at**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **user_id**
  + Type: `int`
- **public**
  + Type: `bool`

### 404

Domain does not exist.
