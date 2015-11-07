define(function(require){
  var Backbone = require('backbone');

	return Backbone.Model.extend({
		fetch: function() {
			return $.when(
				this.get('posts').fetch(),
				this.get('whiskeys').fetch()
      );
		}
	});
})
