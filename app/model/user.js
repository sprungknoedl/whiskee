var UserModel = Backbone.Model.extend({});

var UserCollection = Backbone.Collection.extend({
	url: '/api/users',
	model: UserModel
});

var FriendCollection = Backbone.Collection.extend({
	url: '/api/principal/friends',
	model: UserModel
});
