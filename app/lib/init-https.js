module.exports = function(app) {

  // Force HTTPS and WWW
  if (app.get('env') === 'production') {
    app.use(function (req, res, next) {
      if (req.headers['x-forwarded-proto'] === 'https' 
        || req.headers['cf-visitor'] === '{"scheme":"https"}') {
        return next();
      }
      res.redirect('https://' + req.headers.host + req.url);
    });
  }

};
