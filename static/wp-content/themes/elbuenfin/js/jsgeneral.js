/* Hover Fade Color */
(function(a){a.fn.hoverFadeColor=function(e){var c={color:null,fadeToSpeed:300,fadeFromSpeed:700};e&&a.extend(c,e);this.each(function(){var d=a(this).css("color"),b=c.color;a(this).hover(function(){hoverElem=this;b==null&&a.browser.msie&&a.browser.version.substr(0,1)<9?setTimeout(function(){b==null&&(b=a(hoverElem).css("color"));a(hoverElem).css("color",d);a(hoverElem).stop().animate({color:b},c.fadeToSpeed)},0):(b==null&&(b=a(hoverElem).css("color")),a(hoverElem).css("color",d),a(hoverElem).stop().animate({color:b},
c.fadeToSpeed))},function(){a(hoverElem).stop().animate({color:d},c.fadeFromSpeed)})});return this}})(jQuery);

$(document).ready(function(){

/* FadeIn OnLoad */
	//$('body').hide().fadeIn(700);

/* Delay */
	jQuery.fn.delay = function(time,func){
		return this.each(function(){
			setTimeout(func,time);
		});
	};

/* Header Animaciones Elementos */
//Animación Menú, Logo, Botones Redes
if($.browser.msie && $.browser.version=="7.0") {
}else{
}

});