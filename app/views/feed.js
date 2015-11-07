define(function(require) {
  var Marionette = require('marionette');

  var ItemView = Marionette.ItemView.extend({
    template: '#feed-item-tpl',
    className: 'item',

    events: {
      'click .action-delete': 'removeAction',
    },

    removeAction: function(e) {
      this.model.destroy();
      this.remove();
    }
  });

  return Marionette.CompositeView.extend({
    template: '#feed-tpl',
    childView: ItemView,
    childViewContainer: '#feed-items',

    initialize: function(options) {
      this.collection = options.posts;
		  this.listenTo(this.collection, 'change sync', this.render);
    }
  });
})
