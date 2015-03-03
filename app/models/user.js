var bcrypt   = require('bcrypt');
var mongoose = require('mongoose');

var userSchema = mongoose.Schema({
  email: String,
  password: String,
  created: {type: Date, default: Date.now}
});

userSchema.virtual('gravatar').get(function() {
  var email = this.email.toLowerCase().trim();
  var hash = require('crypto').createHash('md5').update(email);
  return hash.digest('hex');
});

var colors = [
  '#16a085', '#27ae60', '#2980b9', '#8e44ad', '#2c3e50', '#f39c12',
  '#c0392b', '#7f8c8d'
];
userSchema.virtual('color').get(function() {
  var index = this.email.charCodeAt(2) % (colors.length - 1);
  return colors[index];
});

userSchema.methods.checkPassword = function(raw_password, done) {
  bcrypt.compare(raw_password, this.password, function(err, success) {
    if (err || !success) return done(err || new Error('No match'), false);
    done(null, true);
  });
};

userSchema.statics.findByEmail = function(email, done) {
  this.findOne({email: email}, done);
};

var User = mongoose.model('User', userSchema);
module.exports = User;
