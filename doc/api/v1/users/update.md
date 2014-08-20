# PUT /users/:id

Update a user.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **id**
  + Required

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

### Response

### 200 (application/json) 