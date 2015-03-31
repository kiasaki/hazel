import Router from 'express'
import {Application} from '../models'

let router = new Router();
export default router;

router.get('/', function(req, res, next) {
  console.log(Application.fetchAll());
  Application.fetchAll()
    .then(function(result) {
      next(result.toJSON());
    })
    .catch(next);
});


