var Backbone = require('backbone');

var Model = Backbone.Model.extend({});

var Collection = Backbone.Collection.extend({
  url: '/api/users',
  model: Model
});

var Friends = Collection.extend({
  url: '/api/principal/friends',
});

module.exports = {
  Model: Model,
  Collection: Collection,
  Friends: Friends
}
