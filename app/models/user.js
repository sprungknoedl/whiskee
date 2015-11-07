define(function(require) {
  var Backbone = require('backbone');

  var Model = Backbone.Model.extend({});

  var Collection = Backbone.Collection.extend({
    url: '/api/users',
    model: Model
  });

  var Friends = Backbone.Collection.extend({
    url: '/api/principal/friends',
    model: Model
  });

  return {
    Model: Model,
    Collection: Collection,
    Friends: Friends
  }
});
