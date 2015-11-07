define(function(require) {
	var Backbone = require('backbone');

	var Model = Backbone.Model.extend({});
	var Collection = Backbone.Collection.extend({
		url: '/api/whiskeys',
		model: Model,

		comparator: function(a, b) {
			var nameA = a.get('distillery') +
				' ' + a.get('age') +
				' ' + a.get('name')

			var nameB = b.get('distillery') +
				' ' + b.get('age') +
				' ' + b.get('name');

			if (nameA < nameB) { return -1; }
			if (nameA > nameB) { return 1; }
			return 0;
		}
	});

	return {
		Model: Model,
		Collection: Collection
	};
})
