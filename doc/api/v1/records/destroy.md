# DELETE /records/:id

Delete a record.

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

Record is deleted.

### 404 

Record does not exist.
