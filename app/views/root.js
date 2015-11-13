var $ = require('jquery'),
    Marionette = require('backbone.marionette');

var View = Marionette.LayoutView.extend({
  el: 'body',
  template: require('./templates/root.html'),
  regions: {
    content:    '#content',
    navigation: '#navigation',
    sidebar:    '#sidebar',
  },

  initialize: function() {
    this.render();
  },

  showModal: function(view, options) {
    var el = $("#modal");
    el.html(view.render().el);
    el.modal(options);

    view.triggerMethod('attach');
  }
});

module.exports = View;
