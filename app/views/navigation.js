var Marionette = require('backbone.marionette');

var View = Marionette.ItemView.extend({
  template: require('./templates/navigation.html'),
  events: {
    'click .login': 'login',
    'click .logout': 'logout'
  },

  initialize: function(options) {
    this.router = options.app.router;
    this.model = options.principal;
    
    this.listenTo(this.model, 'change sync', this.render);
  },

  login: function(e) {
    e.preventDefault();
    this.router.login();
  },

  logout: function(e) {
    e.preventDefault();
    this.router.logout();
  }
});

module.exports = View;
