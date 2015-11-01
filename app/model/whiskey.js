var WhiskeyModel = Backbone.Model.extend({
	default: {
		id: null,
		distillery: null,
		name: null,
		type: null,
		age: null,
		abv: null,
		size: null,
	},
});

var WhiskeyCollection = Backbone.Collection.extend({
	url: '/api/whiskeys',
	model: WhiskeyModel,
});
