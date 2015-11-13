window.$ = window.jQuery = require("jquery");
var Bootstrap = require('bootstrap'),
    Backbone = require('backbone'),
    Marionette = require('backbone.marionette'),

    Router = require('./router'),
    Security = require('./security'),
    Post = require('./models/post'),
    Whiskey = require('./models/whiskey'),
    User = require('./models/user'),
    Principal = require('./models/principal'),

    RootView = require('./views/root'),
    Navigation = require('./views/navigation'),
    HomeView = require('./views/home'),
    FriendsView = require('./views/friends'),
    PrincipalSidebar = require('./views/principal-sidebar');

var Application = Marionette.Application.extend({
  initialize: function() {
    this.principal = new Principal.Model();

    var injector = {
      app: this,
      principal: this.principal,
    };

    this.router = new Router(injector);
    this.view = new RootView(injector);
    this.security = new Security(injector);
    this.view.showChildView('navigation', new Navigation(injector));

    this.models = {};
    this.models.posts = new Post.Collection();
    this.models.whiskeys = new Whiskey.Collection();
    this.models.users = new User.Collection();
  },

  showIndex: function() {
    var injector = {
      app: this,
      principal: this.principal,
      posts: this.models.posts,
      whiskeys: this.models.whiskeys
    };

    this.view.showChildView('sidebar', new PrincipalSidebar(injector));
    this.view.showChildView('content', new HomeView(injector));
  },

  showFriends: function() {
    var injector = {
      app: this,
      principal: this.principal,
      users: this.models.users
    };

    this.view.showChildView('sidebar', new PrincipalSidebar(injector));
    this.view.showChildView('content', new FriendsView(injector));
  }
});

$(function(){
    window.App = new Application();
    Backbone.history.start();
});
