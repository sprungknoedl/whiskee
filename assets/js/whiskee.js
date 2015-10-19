$('.dropdown').dropdown();

$('#add-post-form').form({
	fields: {
		whiskey: 'empty',
		body: 'empty'
	}
});

$('#add-whiskey-btn').click(function() {
	$('#whiskey-modal')
		.modal({onApprove: function() {
			$(this).find('form').submit();
		}})
		.modal('show');
})

$('.ui.search')
  .search({
    apiSettings: {
      url: '/search?q={query}'
    },
    fields: {
      results: 'items',
      title: 'EMail',
			description: 'ID'
    },
    minCharacters : 3
  })
;
