import express from 'express'
import bodyParser from 'body-parser'

import corsMiddleware from './lib/middlewares/cors'
import notFoundMiddleware from './lib/middlewares/not-found'
import responseMiddleware from './lib/middlewares/response'
import requestLoggerMiddleware from './lib/middlewares/request-logger'

import routes from './routes'

let app = express()
export default app

// Setup database
import './lib/init-bookshelf'

// Setup app
app.set('port', process.env.PORT || 8000)
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({extended: true}))
app.use(corsMiddleware)
app.use(requestLoggerMiddleware)

// Register routes
app.use('/', routes)

// Setup error handler and response handler
app.use(notFoundMiddleware)
app.use(responseMiddleware)
