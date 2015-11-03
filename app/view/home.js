var WhiskeyForm = Marionette.ItemView.extend({
	template: '#whiskey-form-tpl',

	events: {
		'submit form':               'submit',
		'click #action-add-whiskey': 'showAddWhiskey'
	},

	initialize: function(options) {
		this.posts = options.posts;
		this.whiskeys = options.whiskeys;
		this.listenTo(this.whiskeys, 'change sync', this.render);
	},

	onRender: function() {
		this.$('.checkbox').checkbox();
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
		this.$('.dimmer').addClass('active');

		var body = this.$('[name=body]').val();
		var whiskey = this.$('[name=whiskey]').val();
		var security = this.$('[name=security]:checked').val();

		this.posts.create({
			body:     body,
			security: security,
			date:     new Date(),
			user:     app.Principal,
			whiskey:  {id: +whiskey},
		}, {
			wait: true,
			success: function() {
				// clear form
				this.$('[name=body]').val('');
				this.$('.dropdown').dropdown('clear');

				// remove dimmer
				this.$('.dimmer').removeClass('active');
			}
		});
	},

	showAddWhiskey: function(e) {
		e.preventDefault();
		app.Root.showChildView('modals', new WhiskeyAddForm({ model: this.whiskeys }));
	}
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

var ModalView = Marionette.ItemView.extend({
	onRender: function() {
		var modal = $('#modal');
		var html = this.$el.html()
		this.$el.html('');

		modal.html(html);
		modal.find('.checkbox').checkbox();
		modal.find('.dropdown').dropdown()

		modal.modal({ onApprove : this.onApprove.bind(this) });
		modal.modal('show');
	}
});

var WhiskeyAddForm = ModalView.extend({
	template: '#views-home-add-whiskey',

	onApprove: function() {
		var data = $('#modal form').serializeObject();
		// convert some string fields to numbers
		data.age = +data.age;
		data.abv = +data.abv;
		data.size = +data.size;

		this.model.create(data, {wait: true});
	}
});

var HomeView = Marionette.LayoutView.extend({
	template: '#home-tpl',
	regions: {
		form:   '#form',
		feed:   '#feed',
	},

	initialize: function(options) {
		this.posts = options.posts;
		this.whiskeys = options.whiskeys;
	},

	onBeforeShow: function() {
		this.showChildView('form', new WhiskeyForm({ posts: this.posts, whiskeys: this.whiskeys }));
		this.showChildView('feed', new FeedView({ model: this.posts, collection: this.posts }));
	},
});
