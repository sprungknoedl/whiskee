var app = app || new Marionette.Application();

var PrincipalModel = Backbone.Model.extend({
	url: '/api/principal',
	defaults: {
		auth: false,
	},

	is: function(user) {
		return (this.get('id') === user.id);
	}
})

var Auth = Marionette.Object.extend({
	initialize: function (options) {
		this.lock = new Auth0Lock('pH4lp8obYaYmu57LpWuDBQAqBx6469N9', 'whiskee.eu.auth0.com');

		// check for authentication hash
		var hash = this.lock.parseHash(window.location.hash);
		if (hash && hash.error) {
			return alert('There was an error: ' + hash.error + '\n' + hash.error_description);
		}

		if (hash && hash.id_token) {
			var callback = (function (err, profile) {
				if (err) { return alert('There was an error geting the profile: ' + err.message); }
				Backbone.$.ajaxSetup({headers: {'Authorization': 'Bearer ' + hash.id_token}});

				app.Principal.set({
					auth:    true,
					id:      profile.user_id,
					name:    profile.name,
					nick:    profile.nick,
					email:   profile.email,
					picture: profile.picture,
					created: profile.created_at
				});

				app.Principal.save();
				localStorage.setItem('whiskee:token', hash.id_token);
				localStorage.setItem('whiskee:principal', JSON.stringify(app.Principal.toJSON()));

				app.Router.navigate(options.route, {trigger: true, replace: true});
				return;

			}).bind(this);
			this.lock.getProfile(hash.id_token, callback);
		}

		// load stored authentication data
		var token = localStorage.getItem('whiskee:token');
		var principal = localStorage.getItem('whiskee:principal');
		if (token && principal) {
			Backbone.$.ajaxSetup({headers: {'Authorization': 'Bearer ' + token}});
			app.Principal.set(JSON.parse(principal));
		}
	},

	logout: function () {
		localStorage.removeItem('whiskee:token');
		localStorage.removeItem('whiskee:principal');
		window.location.reload();
	},

	login: function() {
		this.lock.show({ scope: 'openid' });
	},
})

var RootView = Marionette.LayoutView.extend({
	el: 'body',
	className: 'ui pusher dimmer page transition',
	template: '#root-tpl',

	regions: {
		nav:     '#nav',
		main:    '#main',
		sidebar: '#sidebar',
		modals:  '#modals-area'
	},

});

var SidebarView = Marionette.ItemView.extend({
	template:  '#views-sidebar',
	className: 'ui overlay sidebar inverted vertical menu',

	initialize: function() {
		this.listenTo(this.model, 'change sync', this.render);
	},

	onRender: function() {
		this.$el.sidebar('setting', {
			dimPage:  false,
			closable: false
		});
		this.$el.sidebar('show');
	}
});

var NavView = Marionette.ItemView.extend({
	template:  '#nav-tpl',
	className: 'ui top fixed inverted blue menu',

	initialize: function() {
		this.listenTo(this.model, 'change sync', this.render);
	},

	events: {
		'click .login': 'login',
	},

	login: function(e) {
		app.Auth.login();
	}
});

var WhiskeyFormModel = Backbone.Model.extend({
	fetch: function() {
		return $.when(
				this.get('posts').fetch(),
				this.get('whiskeys').fetch());
	}
});

$(function() {
	app.Controller = {
		home: function() {
			var posts = new PostCollection(); posts.fetch();
			var whiskeys = new WhiskeyCollection(); whiskeys.fetch();

			app.Root.showChildView('main', new HomeView({ posts: posts, whiskeys: whiskeys }));
		},

		users: function() {
			var users = new UserCollection(); users.fetch();
			app.Root.showChildView('main', new UsersView({ collection: users }));
		}
	};

	app.Router = new Marionette.AppRouter({
		controller: app.Controller,
		appRoutes: {
			'':      'home',
			'users': 'users'
		},
	});

	app.Root = new RootView();
	app.Principal = new PrincipalModel();
	app.VisibleUser = app.Principal;

	app.Auth = new Auth({route: ''});

	app.Root.render();
	app.Root.showChildView('nav', new NavView( {model: app.Principal} ));
	app.Root.showChildView('sidebar', new SidebarView( {model: app.VisibleUser} ));

	app.start();
	Backbone.history.start();
});
