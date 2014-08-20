# DELETE /users/:id

Delete a user.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **id**
  + Required

## Response

### 204

User is deleted.

### 404 

User does not exist.
