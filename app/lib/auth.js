var bcrypt   = require('bcrypt');
var Models   = require('../models');

var COOKIE_KEY = exports.COOKIE_KEY = 'uid';

exports.checkForUser = function(req, res, next) {
  res.locals.authenticated = false;

  if (!(COOKIE_KEY in req.signedCookies)) {
    return next();
  }
  var uid = req.signedCookies[COOKIE_KEY];
  Models.User.findById(uid, function(err, user) {
    if (err) return next(err);

    if (user === null) {
      return next();
    }

    res.locals.authenticated = true;
    res.locals.user = req.currentUser = user;
    res.locals.user.gravatar = req.currentUser.gravatar =
      require('crypto').createHash('md5').update(user.email.toLowerCase().trim()).digest('hex');

    next();
  });
};

exports.requireAuth = function(req, res, next) {
  if (req.path.indexOf('/a/') === 0) return next();
  if (!res.locals.authenticated) return res.redirect('/signin');
  next();
};

exports.login = function(res, user_id) {
  res.cookie(COOKIE_KEY, user_id.toString(), {signed: true});
};

exports.logout = function(res) {
  res.clearCookie(COOKIE_KEY, {signed: true});
};

exports.tryLogin = function(res, username, password, callback) {
  Models.User.findByEmail(username, function(err, user) {
    if (err) return callback(err, false);
    if (user === null) return callback(new Error('Can\'t find user'), false);

    bcrypt.compare(password, user.password, function(err, result) {
      if (err) return callback(err, false);
      if (!result) return callback(new Error('Bad password'), false);

      exports.login(res, user._id);
      callback(null, true);
    });
  });
};
