var Marionette = require('backbone.marionette'),
    AddWhiskeyModal = require('./modals/add-whiskey'),
    moment = require('moment');
    require('bootstrap-select');

var ChildView = Marionette.ItemView.extend({
  className: 'media',
  template: require('./templates/content/home-post.html'),
  templateHelpers: {
    moment: moment,
    security: null, // initialized later
    model: function() { return this; }
  },

  events: {
    'click a.edit': 'onEdit',
    'click a.delete': 'onDelete',
  },

  initialize: function(options) {
    this.templateHelpers.security = options.security;
  },

  onEdit: function(e) { e.preventDefault(); },
  onDelete: function(e) {
    e.preventDefault();
    this.model.destroy();
    this.remove();
  }
});

var View = Marionette.CompositeView.extend({
  template: require('./templates/content/home.html'),
  childView: ChildView,
  childViewContainer: '#child-view',

  events: {
    'submit form': 'onSubmit',
    'click #add-whiskey': 'onAddWhiskey'
  },

  initialize: function(options) {
    this.posts = options.posts;
    this.whiskeys = options.whiskeys;
    this.parent = options.app.view;
    this.collection = options.posts;

    this.refresh();

    this.listenTo(this.posts, 'sync change', this.render);
    this.listenTo(this.whiskeys, 'sync change', this.render);
    this.childViewOptions = { security: options.app.security }
    options.app.vent.on('login', this.refresh.bind(this));
  },

  refresh: function() {
    this.posts.fetch();
    this.whiskeys.fetch();

  },

  serializeData: function() { return { whiskeys: this.whiskeys.toJSON() }; },
  onRender: function() { this.$('select').selectpicker(); },

  onAddWhiskey: function(e) {
    e.preventDefault();
    var modal = new AddWhiskeyModal();
    this.parent.showModal(modal);
  },

  onSubmit: function(e) {
    e.preventDefault();
    this.posts.create({
      whiskey:  {id: +this.$('[name=whiskey]').val()},
      body:     this.$('[name=body]').val(),
      security: this.$('[name=security]:checked').val()
    }, {wait: true});

    // clear form
    this.$('form').trigger('reset');
    this.$('select').selectpicker('refresh');
  }
});

module.exports = View;
