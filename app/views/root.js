define(function(require) {
  var Marionette = require('marionette'),
      NavigationView = require('views/navigation'),
      SidebarView = require('views/sidebar'),
      App = require('app');

	return Marionette.LayoutView.extend({
		el: 'body',
		className: 'ui pusher dimmer page transition',
		template: '#root-tpl',

		regions: {
			nav:     '#nav',
			main:    '#main',
			sidebar: '#sidebar',
			modals:  '#modals-area'
		},

    onRender: function() {
      this.showChildView('nav', new NavigationView());
      this.showChildView('sidebar', new SidebarView());
    }
	});
})
