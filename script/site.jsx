if(!localStorage['lastPage'] || localStorage['lastPage'] != window.location.href) {
	localStorage['lastPage'] = window.location.href;
	$('body').fadeIn('slow');
} else {
	$('body').show();
}
