$(document).ready(function(){
	var menuVisible = false;

	$(".navigation-selection-item > span").on('click', function(){
		$('.navigation-item > a').css('display', menuVisible ? 'none' : 'block');
		menuVisible = !menuVisible
		$(this).toggleClass('selected')
	});
});