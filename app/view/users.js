var SingleUserView = Marionette.ItemView.extend({
	template:  '#single-user-tpl',
	className: 'ui item',

	events: {
		'click .action-friend': 'addFriend'
	},
	
	addFriend: function(e) {
		console.log("added friend", this.model.get('name'));
		new FriendCollection().create({id: this.model.get('id')});
	}
});

var UsersView = Marionette.CompositeView.extend({
	template:           '#users-tpl',
	childView:          SingleUserView,
	childViewContainer: '#child-view',
});
