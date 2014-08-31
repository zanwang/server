# POST /tokens/facebook

Log in with facebook.

## Request

### Body

- **user_id**
  + Required
  + Type: `string`
- **access_token**
  + Required
  + Type: `string`

## Response

### 201 (application/json)

- **user_id**
  + Type: `int`
- **key**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **expired_at**
  + Type: `string`

### 400

Email has been taken.