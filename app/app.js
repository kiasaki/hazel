try {
  var config = require('./config');
  for (var k in config) {
    process.env[k] = config[k];
  }
} catch (e) {}

var express      = require('express');
var path         = require('path');
var favicon      = require('serve-favicon');
var logger       = require('morgan');
var cookieParser = require('cookie-parser');
var bodyParser   = require('body-parser');
var debug        = require('debug')('hazel');
var auth         = require('./lib/auth');

var app = express();

app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

require('./lib/init-https')(app);
app.use(favicon(path.join(__dirname, 'public', 'favicon.ico')));
app.use(logger('dev'));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));
app.use(cookieParser(process.env.SECRET || 'keyboardcat'));
app.use(express.static(path.join(__dirname, 'public')));

var routes      = require('./routes/index');
var usersRoutes = require('./routes/users');
var authRoutes = require('./routes/auth');

// Public Routes
app.use(auth.checkForUser);
app.use('/', authRoutes);

// Auth protected routes
app.use(auth.requireAuth);
app.use('/', routes);
app.use('/users', usersRoutes);

// Handle 404s and errors sent down the line
require('./lib/init-errors')(app);

if (!module.parent) {
  require('./lib/db').setup(function() {
    app.listen(process.env.PORT || 8080, '0.0.0.0', function() {
      debug('Server listening on port ' + this.address().port);
    });
  });
}

module.exports = app;
