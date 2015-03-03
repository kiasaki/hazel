var mongoose = require('mongoose');

exports.setup = function(callback) {
  var dbUrl = process.env.MONGOLAB_URI || 'mongodb://localhost/hazel';
  mongoose.connect(dbUrl);

  var db = mongoose.connection;
  db.on('error', console.error.bind(console, 'connection error:'));
  db.once('open', callback);
};
