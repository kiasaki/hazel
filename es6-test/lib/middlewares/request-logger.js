export default function requestLoggerMiddleware(req, res, next) {
  let date = new Date().toISOString().split('.')[0]
  let url = req.originalUrl || req.url

  console.log(`${date} ${req.method} ${url}`)
  next()
}
