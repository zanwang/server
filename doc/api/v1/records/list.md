# GET /domains/:domain_id/records

List records of a domain.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **domain_id**
  + Required

## Response

### 200 (application/json)

Returns a list of records.

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

Domain does not exist. 