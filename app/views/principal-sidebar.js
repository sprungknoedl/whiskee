var Marionette = require('backbone.marionette'),
    moment = require('moment');

var View = Marionette.ItemView.extend({
  template: require('./templates/sidebars/principal.html'),
  templateHelpers: {
    moment: moment,
  },
  
  initialize: function(options) {
    this.model = options.principal;
    this.listenTo(this.model, 'change sync', this.render);
  }
});

module.exports = View;
