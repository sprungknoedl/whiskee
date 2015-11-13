var Backbone = require('backbone');

var Model = Backbone.Model.extend({
  url: '/api/principal',
  defaults: {
    auth: false,
  }
});

module.exports = {
  Model: Model
};
