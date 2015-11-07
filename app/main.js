require.config({
  shim: {
    underscore: {
      exports: '_'
    },
    backbone: {
      deps: [ 'underscore', 'jquery' ],
      exports: 'Backbone'
    },
    marionette: {
      deps: [ 'backbone' ],
      exports: 'Marionette'
    },
    'lib/jquery-jsonform': ['jquery'],
    'lib/jquery-transloadit': ['jquery'],
    'lib/semantic': ['jquery']
  },
  paths: {
    jquery: 'lib/jquery',
    underscore: 'lib/underscore',
    backbone: 'lib/backbone',
    marionette: 'lib/marionette',
    'auth0': 'lib/auth0'
  }
});

define(function(require) {
  var App = require('app'),
      Marionette = require('marionette'),
      RootView = require('views/root'),
      HomeView = require('views/home'),
      UsersView = require('views/users'),
      PrincipalModel = require('models/principal');

  App.Root = new RootView();
  App.Principal = new PrincipalModel();

  App.Controller = {
    home: function() {
      App.Root.showChildView('main', new HomeView());
    },

    users: function() {
      App.Root.showChildView('main', new UsersView());
    }
  };

  App.Router = new Marionette.AppRouter({
    controller: App.Controller,
    appRoutes: {
      ''        : 'home',
      'users'   : 'users'
    }
  });

  App.onStart = function() {
    App.Root.render();
	  Backbone.history.start();
  }

  App.start();
});
