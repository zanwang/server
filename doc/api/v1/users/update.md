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
  + Optional
  + Type: `string`
- **password**
  + Optional
  + Type: `string`
  + Length: 6~50
- **old_password**
  + Optional
  + Type: `string`
  + Length: 6~50
- **email**
  + Optional
  + Type: `string`

### Response

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

### 403

Password is wrong. 