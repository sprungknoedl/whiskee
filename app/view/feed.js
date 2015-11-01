var WhiskeyForm = Marionette.ItemView.extend({
	template: '#whiskey-form-tpl',

	events: {
		'submit form':  'submit',
	},

	initialize: function(options) {
		this.posts = options.posts;
		this.whiskeys = options.whiskeys;
		this.listenTo(this.whiskeys, 'change sync', this.render);
	},

	onRender: function() { 
		this.$('.dropdown').dropdown();
	},

	serializeData: function() {
		return {
			posts:    this.posts.toJSON(),
			whiskeys: this.whiskeys.toJSON()
		};
	},

	submit: function(e) {
		e.preventDefault();
		var body = this.$('[name=body]').val();
		var whiskey = this.$('[name=whiskey]').val();

		// clear form
		this.$('[name=body]').val('');
		this.$('.dropdown').dropdown('clear');

		var post = this.posts.create({
			body:     body,
			security: 'public',
			date:     new Date(),
			user:     app.Principal,
			whiskey:  {id: +whiskey},
		}, { wait: true });
	},
});

var FeedItemView = Marionette.ItemView.extend({
	template:  '#feed-item-tpl',
	className: 'item',

	events: {
		'click .action-delete': 'removeAction',
	},

	removeAction: function(e) {
		this.model.destroy();
		this.remove();
	}
});

var FeedView = Marionette.CompositeView.extend({
	template:           '#feed-tpl',
	childView:          FeedItemView,
	childViewContainer: '#feed-items',
});

var HomeView = Marionette.LayoutView.extend({
	template: '#home-tpl',
	regions: {
		form:   '#form',
		feed:   '#feed',
	},

	events: {
		'click .login': 'login',
	},

	initialize: function(options) {
		this.posts = options.posts;
		this.whiskeys = options.whiskeys;
	},

	onBeforeShow: function() {
		this.showChildView('form', new WhiskeyForm({ posts: this.posts, whiskeys: this.whiskeys }));
		this.showChildView('feed', new FeedView({ model: this.posts, collection: this.posts }));
	},

	login: function(e) {
		app.Auth.login();
	}
});
