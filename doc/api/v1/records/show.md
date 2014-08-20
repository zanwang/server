# GET /records/:id

Get a record.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **id**
  + Required

## Response

### 200 (application/json)

- **id**
  + Type: `int`
- **name**
  + Type: `string`
- **type**
  + Type: `string`
  + Possible values: `A`, `CNAME`, `MX`, `TXT`, `SPF`, `AAAA`, `NS`, `LOC`
- **value**
  + Type: `string`
- **created_at**
  + Type: `string`
- **updated_at**
  + Type: `string`
- **domain_id**
  + Type: `int`
- **ttl**
  + Type: `uint`
- **priority**
  + Type: `uint`

### 404

Record does not exist.
