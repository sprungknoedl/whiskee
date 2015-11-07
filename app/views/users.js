define(function(require) {
  var Marionette = require('marionette'),
      Users = require('models/user');

  var ChildView = Marionette.ItemView.extend({
    template: '#single-user-tpl',
    className: 'ui item',

    events: {
      'click .action-friend': 'addFriend'
    },

    addFriend: function(e) {
      console.log("added friend", this.model.get('name'));
      new Users.Friends().create({
        id: this.model.get('id')
      });
    }
  });

  return Marionette.CompositeView.extend({
    template: '#users-tpl',
    childView: ChildView,
    childViewContainer: '#child-view',

    initialize: function() {
      this.collection = new Users.Collection();
      this.collection.fetch();
    }
  });
})
