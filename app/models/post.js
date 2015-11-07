define(function(require) {
	var Backbone = require('backbone');

	var Model = Backbone.Model.extend({});
	var Collection = Backbone.Collection.extend({
		url: '/api/posts',
		model: Model,
		comparator: function(a, b) {
			if (a.get('date') > b.get('date')) { return -1; }
			if (a.get('date') < b.get('date')) { return 1; }
			return 0;
		}
	});

	return {
		Model: Model,
		Collection: Collection,
	};
});
