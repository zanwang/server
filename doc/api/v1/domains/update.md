# PUT /domains/:id

Update a domain.

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
  + Maximum length: 63
  + Must be a valid domain name

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

### 400

Domain name has been taken.

### 404

Domain does not exist.
