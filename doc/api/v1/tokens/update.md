# PUT /tokens

Update a token.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

## Response

### 200 (application/json)

- **user_id**
  + Type: `int`
- **key**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **expired_at**
  + Type: `string`
