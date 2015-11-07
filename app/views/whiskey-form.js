define(function(require) {
  var App = require('app'),
      Marionette = require('marionette');
  require('lib/jquery-jsonform');
  require('lib/jquery-transloadit');

  var ModalView = Marionette.ItemView.extend({
    onRender: function() {
      var modal = $('#modal');
      var html = this.$el.html()
      this.$el.html('');

      modal.html(html);
      modal.find('.checkbox').checkbox();
      modal.find('.dropdown').dropdown()

      modal.modal({
        onApprove: this.onApprove.bind(this)
      });
      modal.modal('show');
    }
  });

  var WhiskeyAddForm = ModalView.extend({
    template: '#views-home-add-whiskey',

    onRender: function() {
      ModalView.prototype.onRender.apply(this);

      $('#modal form').transloadit({
        wait: true,
        autoSubmit: false,
        triggerUploadOnFileSelection: true,

        params: {
          auth: { key: "61eec27083fd11e5900d9dd1b00d757c" },
          template_id: "fe186c20847d11e596671710a5660bd7"
        }
      });
    },

    onApprove: function() {
      var data = $('#modal form').serializeObject();

      // store pictures
      if (data.transloadit) {
        var transloadit = JSON.parse(data.transloadit);
        data.picture = transloadit.results[':original'][0].ssl_url;
        data.thumb = transloadit.results['resize'][0].ssl_url;
      } else {
        data.picture = '/static/img/default.jpg';
        data.thumb = '/static/img/default-thumb.jpg';
      }

      // modify data object before sending
      data.age = +data.age;
      data.abv = +data.abv;
      data.size = +data.size;
      data.transloadit = null;

      this.model.create(data, {
        wait: true
      });
    }
  });

  return Marionette.ItemView.extend({
    template: '#whiskey-form-tpl',

    events: {
      'submit form': 'submit',
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
        posts: this.posts.toJSON(),
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
        body: body,
        security: security,
        date: new Date(),
        user: App.Principal,
        whiskey: { id: +whiskey },
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
      App.Root.showChildView('modals', new WhiskeyAddForm({
        model: this.whiskeys
      }));
    }
  });
});
