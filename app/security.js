var Marionette = require('backbone.marionette');

var Security = Marionette.Object.extend({
  initialize: function(options) {
    this.principal = options.principal;
  },

  isAllowed: function(action, target) {
    if (!this.principal.get('auth')) { return false; }

    var principal = this.principal.get('id');
    switch(action) {
      case 'post:edit':
        return principal == target.user.id;
      case 'post:delete':
        return principal == target.user.id;
    }

    console.log('can\'t authorize unknown action', action);
    return false;
  }
});

module.exports = Security;
