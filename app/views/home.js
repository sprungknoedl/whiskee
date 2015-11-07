define(function(require) {
  var Marionette = require('marionette'),
      Posts = require('models/post'),
      Whiskeys = require('models/whiskey'),
      WhiskeyForm = require('views/whiskey-form'),
      FeedView = require('views/feed');

  return Marionette.LayoutView.extend({
    template: '#home-tpl',
    regions: {
      form: '#form',
      feed: '#feed',
    },

    initialize: function() {
      this.posts = new Posts.Collection();
      this.whiskeys = new Whiskeys.Collection();

      this.posts.fetch();
      this.whiskeys.fetch();
    },

    onBeforeShow: function() {
      this.showChildView('form', new WhiskeyForm({
        posts: this.posts,
        whiskeys: this.whiskeys
      }));

      this.showChildView('feed', new FeedView({
        posts: this.posts,
      }));
    },
  });
});
