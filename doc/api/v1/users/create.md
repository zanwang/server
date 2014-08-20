# POST /users

Create a user.

## Request

### Body

- **name**
  + Required
  + Type: `string`
- **password**
  + Required
  + Type: `string`
  + Length: 6~50
- **email**
  + Required
  + Type: `string`

## Response

### 201 (application/json)

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

### 400 (application/json)

Email has been taken.