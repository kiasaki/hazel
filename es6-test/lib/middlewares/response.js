import {ErrorNotFound} from '../errors'
import {STATUS_CODES} from 'http'

export default function responseMiddleware(err, req, res, next) {
  if (err instanceof Error) {
    respondWithError(res, err)
  } else {
    res.json(err);
  }
}

function respondWithError(res, err) {
  let statusCode = 500

  if (err.statusCode && res.statusCode in STATUS_CODES) {
    statusCode = err.statusCode
  }

  res.status(statusCode)
  res.json({
    error: {
      type: err.type || '',
      message: err.message,
      code: statusCode
    }
  })
}
