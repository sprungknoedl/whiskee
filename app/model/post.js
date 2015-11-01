var PostModel = Backbone.Model.extend({
	urlRoot: '/api/posts'
});

var PostCollection = Backbone.Collection.extend({
	url: '/api/posts',
	model: PostModel
});
