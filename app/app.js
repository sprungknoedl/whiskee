var app = app || new Marionette.Application();

var PrincipalModel = Backbone.Model.extend({
	url: '/api/principal'
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

				app.Principal = new PrincipalModel({
					id:      profile.user_id,
					name:    profile.name,
					nick:    profile.nick,
					email:   profile.email,
					picture: profile.picture
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
			app.Principal = new PrincipalModel(JSON.parse(principal));
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
	el: '#body',
	template: '#root-tpl',

	regions: {
		sidebar: '#sidebar',
		main:    '#main',
	}
});

var WhiskeyFormModel = Backbone.Model.extend({
	fetch: function() {
		return $.when(
				this.get('posts').fetch(), 
				this.get('whiskeys').fetch());
	}
});

app.Controller = {
	home: function() {
		var posts = new PostCollection(); posts.fetch()
		var whiskeys = new WhiskeyCollection(); whiskeys.fetch()

		app.Root.showChildView('main', new HomeView({ posts: posts, whiskeys: whiskeys }));
	},
}

app.Router = new Marionette.AppRouter({
	controller: app.Controller,
	appRoutes: {
		'': 'home',
	},
});

app.Root = new RootView();
app.Auth = new Auth({route: ''});

$(function() {
	app.start();
	app.Root.render();

	Backbone.history.start();
});

