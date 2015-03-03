var express = require('express');
var auth    = require('../lib/auth');
var Models  = require('../models');
var router  = express.Router();

router.get('/signin', function(req, res) {
  res.render('auth/signin', {});
});

router.post('/signin', function(req, res) {
  var email = req.param('email');
  var password = req.param('password');

  if (!email || !password) {
    res.render('auth/signin', {
      email: email,
      error: 'Email and password are necessary'
    });
    return;
  }

  auth.tryLogin(res, email, password, function(err, worked) {
    if (err || !worked) {
      res.render('auth/signin', {
        email: email,
        error: 'Email or password invalid'
      });
      return;
    }

    res.redirect('/');
  });
});

router.get('/signout', function(req, res) {
  auth.logout(res);
  res.redirect('/');
});

router.get('/forgot', function(req, res) {
  res.render('auth/forgot', {});
});

module.exports = router;
