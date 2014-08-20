# POST /domains/:domain_id/records

Create a record.

## Request

### Headers

- **Authorization**
  + Required
  + Format: `token <token>`

### Params

- **domain_id**
  + Required

### Body

- **name**
  + Required
  + Type: `string`
  + Maximum length: 63
  + Must be a valid domain name
- **type**
  + Required
  + Type: `string`
  + Possible values: `A`, `CNAME`, `MX`, `TXT`, `SPF`, `AAAA`, `NS`, `LOC`
- **value**
  + Required
  + Type: `string`
- **ttl**
  + Required
  + Type: `uint`
  + The value must be 1 (Automatic) or 300~86400 (seconds)
- **priority**
  + Optional
  + Type: `uint`

## Response

### 201 (application/json)

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
