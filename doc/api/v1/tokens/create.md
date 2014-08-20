# POST /tokens

Create a token.

## Request

### Body

- **email**
  + Required
  + Type: `string`
- **password**
  + Required
  + Type: `string`
  + Length: 6~50

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

User does not exist.

### 401

Password is wrong.