define(function(require) {
  var Backbone = require('backbone');

  return Backbone.Model.extend({
    url: '/api/principal',
    defaults: {
      auth: false,
    },

    is: function(user) {
      return (this.get('id') === user.id);
    }
  });
})
