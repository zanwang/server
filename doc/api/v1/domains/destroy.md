# DELETE /domains/:id

Delete a domain.

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

Domain is deleted.

### 404

Domain does not exist.
