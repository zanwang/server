# POST /users/:user_id/domains

Create a domain.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **user_id**
  + Required

### Body

- **name**
  + Required
  + Type: `string`
  + Maximum length: 63
  + Must be a valid domain name

## Response

### 201 (application/json)

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

### 400

Domain name has been taken.

### 404

User does not exist.
