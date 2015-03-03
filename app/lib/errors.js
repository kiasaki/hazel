module.exports = function(code, message) {
  var err = new Error(message);
  err.status = code;
  return err;
};
