var Backbone = require('backbone'),
    Marionette = require('backbone.marionette'),
    Auth0Lock = require('auth0-lock'),
    Security = require('./security');

var Router = Marionette.AppRouter.extend({
  appRoutes: {
    '':        'showIndex',
    'friends': 'showFriends'
  },

  initialize: function(options) {
    // store required state
    this.controller = options.app;
    this.principal = options.principal;
		this.lock = new Auth0Lock(
      'pH4lp8obYaYmu57LpWuDBQAqBx6469N9',
      'whiskee.eu.auth0.com'
    );

    // check for authentication information
		var hash = this.lock.parseHash(window.location.hash);
		if (hash && hash.error) {
			return alert('There was an error: ' + hash.error + '\n' + hash.error_description);
		}

    // load profile from auth0 and authenticate user
		if (hash && hash.id_token) {
    	var callback = function (err, profile) {
    		if (err) { return alert('There was an error geting the profile: ' + err.message); }
    		var principal = {
    			auth:    true,
    			id:      profile.user_id,
    			name:    profile.name,
    			nick:    profile.nick,
    			email:   profile.email,
    			picture: profile.picture,
    			created: profile.created_at
    		};

        this.authenticate(hash.id_token, principal, true);
    	};

			this.lock.getProfile(hash.id_token, callback.bind(this));
      window.location.hash = ''; // TODO: find a way to do this via backbone
      return;
		}

		// load stored authentication data
		var token = localStorage.getItem('whiskee:token');
		var principal = localStorage.getItem('whiskee:principal');
		if (token && principal) { this.authenticate(token, JSON.parse(principal)); }
  },

  authenticate: function(token, principal, register) {
		$.ajaxSetup({headers: {'Authorization': 'Bearer ' + token}});
		this.principal.set(principal);
    this.controller.vent.trigger('login', this.principal);

    if (register) {
      this.principal.save();
  		localStorage.setItem('whiskee:token', token);
  		localStorage.setItem('whiskee:principal', JSON.stringify(principal));
    }
  },

	login: function() {
		this.lock.show({ scope: 'openid' });
	},

	logout: function() {
		localStorage.removeItem('whiskee:token');
		localStorage.removeItem('whiskee:principal');
		window.location.reload();
	}
});

module.exports = Router;
