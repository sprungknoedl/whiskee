var PostModel = Backbone.Model.extend({});

var PostCollection = Backbone.Collection.extend({
	url: '/api/posts',
	model: PostModel,
	comparator: function(a, b) {
		if (a.get('date') > b.get('date')) { return -1; }
		if (a.get('date') < b.get('date')) { return 1; }
		return 0;
	}
});
