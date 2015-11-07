define(function(require) {
  var $ = require('jquery'),
      Marionette = require('marionette'),
      App = require('app');

	var Auth = Marionette.Object.extend({
		initialize: function () {
			this.lock = new Auth0Lock('pH4lp8obYaYmu57LpWuDBQAqBx6469N9', 'whiskee.eu.auth0.com');

			// check for authentication hash
      App.on('before:start', this.checkAuth.bind(this));
		},

    checkAuth: function () {
			var hash = this.lock.parseHash(window.location.hash);
			if (hash && hash.error) {
				return alert('There was an error: ' + hash.error + '\n' + hash.error_description);
			}

			if (hash && hash.id_token) {
				var callback = (function (err, profile) {
					if (err) { return alert('There was an error geting the profile: ' + err.message); }
					$.ajaxSetup({headers: {'Authorization': 'Bearer ' + hash.id_token}});

					App.Principal.set({
						auth:    true,
						id:      profile.user_id,
						name:    profile.name,
						nick:    profile.nick,
						email:   profile.email,
						picture: profile.picture,
						created: profile.created_at
					});

					App.Principal.save();
					localStorage.setItem('whiskee:token', hash.id_token);
					localStorage.setItem('whiskee:principal', JSON.stringify(App.Principal.toJSON()));

					App.Router.navigate('', {trigger: true, replace: true});
					return;

				}).bind(this);
				this.lock.getProfile(hash.id_token, callback);
			}

			// load stored authentication data
			var token = localStorage.getItem('whiskee:token');
			var principal = localStorage.getItem('whiskee:principal');
			if (token && principal) {
        console.log('logging in user');
				$.ajaxSetup({headers: {'Authorization': 'Bearer ' + token}});
				App.Principal.set(JSON.parse(principal));
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
	});

  return new Auth();
})
