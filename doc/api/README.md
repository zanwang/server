# API Overview

Base URL: `https://maji.moe/api/<version>`

- [v1 (latest)](v1/README.md)

## Synopsis

- All responses are JSON.
- You can use form (`application/x-www-form-urlencoded`) or JSON (`application/json`) in requests

## Authorization

The API is stateless. You have to request a token first. For example:

``` js
POST /api/v1/tokens
{
    "email": "abc@maji.moe",
    "password": "123456"
}
```

And put the token in header when needed.

```
Authorization: token <token>
```

## Error Response

### Example

``` js
{
    "code": 111,
    "field": "name",
    "message": "Name is required"
}
```

### HTTP Status Code

- 400 Bad request - Something must be wrong in your request.
- 401 Unauthorized - You need to send your request again with a valid token.
- 403 Forbidden - You are forbidden to access the content.
- 404 Not found - The data you want to access is deleted or doesn't exist.
- 500 Server error - Server may be down.

### Error Code

- 110: Unknown
- 111: Required
- 112: Unsupported content type
- 113: JSON parsing error
- 114: Wrong data type
- 120: Email
- 121: URL
- 122: Alpha
- 123: Alphanumeric
- 124: Numeric
- 125: Hexadecimal
- 126: Hex color
- 127: Lowercase
- 128: Uppercase
- 129: Int
- 130: Float
- 131: Divisble
- 132: Length
- 133: Min length
- 134: Max length
- 135: UUID
- 136: Credit card
- 137: ISBN
- 138: JSON
- 139: Multibyte
- 140: ASCII
- 141: Full width
- 142: Half width
- 143: Variable width
- 144: Base64
- 145: IP
- 146: IPv4
- 147: IPv6
- 148: MAC
- 149: Min value
- 150: Max value
- 151: Range
- 152: Domain name
- 153: Domain
- 210: User hasn't been activated
- 211: User has been activated
- 212: Email has been taken
- 213: Domain has been taken
- 214: Wrong password
- 215: Password hasn't been set
- 216: Wrong record type
- 217: Token has been expired
- 218: User doesn't exist
- 219: Token is required
- 220: Domain doesn't exist
- 221: Forbidden to access the domain
- 222: Domain doesn't exist
- 223: Forbidden to access the record
- 224: Domain has been reserved
- 225: Forbidden to access the user
- 226: Domain can't be renew
- 500: Server error