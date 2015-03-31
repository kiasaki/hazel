require('babel/register');

var app = require('./app');

app.listen(app.get('port'));
console.log('Hazel API listening on port %s', app.get('port'));
