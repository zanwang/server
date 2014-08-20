# API Overview

- Base URL: `http://maji.moe/api/v1`

## Error Response

### 1XX - Format error

- 110: Required
- 111: Minimum length
- 112: Maximum length
- 113: Length
- 114: Email
- 115: Minimum value
- 116: Maximum value
- 117: Within a range
- 118: Out of a range
- 119: In an array
- 120: Not in an array
- 121: IP (v4 and v6)
- 123: Domain
- 124: Domain name

### 2XX - Custom Error

- 210: User has not been activated
- 211: Email has been taken
- 212: Domain name has been taken
- 213: User does not exist
- 214: Password is wrong
- 215: User has been activated
