$(document).ready(function(){

/* FadeIn OnLoad */
	$('body').hide().fadeIn(600);
	
/* Slider Principal */
if($.browser.msie) {
	$('#slider').nivoSlider({
		effect: 'fade',
		directionNavHide: false, 
		controlNav: false,
		pauseTime: 5000
	});
}else{
	$('#slider').nivoSlider({
		effect: 'boxRandom',
		directionNavHide: false,
		controlNav: false,
		pauseTime: 5000
	});
}

/* Carousel Entidades Participantes */
	$("div.partC").carousel({ 
		dispItems: 3,
		autoSlide: true,
        autoSlideInterval: 3000,
        delayAutoSlide: 1500,
		animSpeed: "slow",
		loop: true
	});

});