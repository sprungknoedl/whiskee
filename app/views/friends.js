var Marionette = require('backbone.marionette'),
    moment = require('moment');

var View = Marionette.ItemView.extend({
  template: require('./templates/content/friends.html'),
  templateHelpers: { moment: moment },

  initialize: function(options) {
    this.model = options.users;
    this.refresh();

    this.listenTo(this.model, 'sync change', this.render);
  },

  refresh: function() { this.model.fetch(); },
  serializeData: function() { return { users: this.model.toJSON() }; }
});

module.exports = View;
