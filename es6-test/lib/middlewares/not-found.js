import {ErrorNotFound} from '../errors'

export default function notFoundMiddleware(req, res, next) {
  next(new ErrorNotFound());
}
