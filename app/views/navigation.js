define(function(require) {
  var App = require('app'),
      Auth = require('controller/auth'),
      Marionette = require('marionette');

	return Marionette.ItemView.extend({
		template:  '#nav-tpl',
		className: 'ui top fixed inverted blue menu',

		events: {
			'click .login': 'login',
			'click .logout': 'logout'
		},

		initialize: function() {
      this.model = App.Principal;
			this.listenTo(this.model, 'change sync', this.render);
		},

		login: function(e) { Auth.login(); },
		logout: function(e) { Auth.logout(); }
	});
})
