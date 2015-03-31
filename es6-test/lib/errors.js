import {STATUS_CODES} from 'http'

export class HTTPError extends Error {
  constructor(code, message) {
    this.type = STATUS_CODES[code]
    this.statusCode = code
    super(message || this.type)
  }
}

export class ErrorBadRequest extends HTTPError {
  constructor(message) {
    super(400, message)
  }
}

export class ErrorUnauthorized extends HTTPError {
  constructor(message) {
    super(401, message)
  }
}

export class ErrorNotFound extends HTTPError {
  constructor(message) {
    super(404, message)
  }
}

export class ErrorInternalServer extends HTTPError {
  constructor(message) {
    super(500, message);
  }
}
