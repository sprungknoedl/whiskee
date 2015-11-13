var Modal = require('./modal');
    require('jquery-transloadit2');

var View = Modal.extend({
  template: require('../templates/modals/add-whiskey.html'),

  onAttach: function() {
    $('#modal select').selectpicker();
    $('#modal form').transloadit({
      wait: true,
      autoSubmit: false,
      triggerUploadOnFileSelection: true,

      params: {
        auth: { key: "61eec27083fd11e5900d9dd1b00d757c" },
        template_id: "fe186c20847d11e596671710a5660bd7"
      }
    });
  }
});

module.exports = View;
