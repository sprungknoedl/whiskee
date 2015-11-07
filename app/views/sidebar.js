define(function(require) {
  var App = require('app'),
      Marionette = require('marionette');
      require('lib/semantic');

	return Marionette.ItemView.extend({
		template:  '#views-sidebar',
		className: 'ui overlay sidebar inverted vertical menu',

		initialize: function() {
      this.model = App.Principal;
			this.listenTo(this.model, 'change sync', this.render);
		},

		onRender: function() {
			this.$el.sidebar('setting', {
				dimPage:  false,
				closable: false
			});
			this.$el.sidebar('show');
		}
	});
})
