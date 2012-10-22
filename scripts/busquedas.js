$(document).ready(function(){       
	getcarrousel();
	pagina = 0;
	var queryVars = getVars();
	var hasVars = false;
	//if(queryVars.hasOwnProperty('word') || queryVars.hasOwnProperty('catMenu') || queryVars.hasOwnProperty('estadoMenu') || queryVars.hasOwnProperty('tipoMenu') )
	if(queryVars.hasOwnProperty('word')){
		hasVars = true;
		$("input[name=word]").val(queryVars.word);
	}
	if(queryVars.hasOwnProperty('catMenu')){
		hasVars = true;
		$("select[name=catMenu]").val(queryVars.catMenu);
	}
	if(queryVars.hasOwnProperty('estadoMenu')){
		hasVars = true;
		$("select[name=estadoMenu]").val(queryVars.estadoMenu);
	}
	if(queryVars.hasOwnProperty('tipoMenu')){
		hasVars = true;
		$("select[name=tipoMenu]").val(queryVars.tipoMenu);
	}        
	search();
	$("#buscarOferta").click(function(){                
		pagina = 0;    
		search();
		return false;
	});
	cargaOfertas = false;
	$(document).scroll(function(){
		if(cargaOfertas){
			if(($(document).scrollTop() + $(window).height()) >= ($(document).height() - 10)){
				search();
			}
		}
	})

	$("#CloserLight").click(function() {
		$('#cuerpo').removeClass('noscroll');
		$('#lightback').addClass("hide"); 
		$('#lightfront').addClass("hide"); 
				  return false;
	});

	$(".lighter").click(function() {
		$('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
		$('#lightback').removeClass("hide"); 
		$('#lightfront').removeClass("hide"); 
				  return false;
	});
	});
	function getVars() {
	var delimiter = "?"; // using '#' here is great for AJAX apps.
	var separator = "&";
	var url = location.href;
	var get_exists = (url.indexOf(delimiter) > -1) ? true : false;
	if (get_exists) {
		var url_get = {};
		var params = url.substring(url.indexOf(delimiter)+1);
		var params_array = params.split(separator);
		for (var i=0; i < params_array.length; i++) {
			var param_name = params_array[i].substring(0,params_array[i].indexOf('='));
			var param_value = params_array[i].substring(params_array[i].indexOf('=')+1);
			url_get[param_name] = param_value;
		}
	  return url_get;
	}
	return false;
	}
	function lighterAjax(){
	$(".lighter").click(function() {
		$('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
		$('#lightback').removeClass("hide"); 
		$('#lightfront').removeClass("hide"); 
	  return false;
	});
	}
	function getcarrousel() {
	$.get("/carr", "", function(response){
		$('#logo1').html(response);
	});
	}
	function search(){  
	var keywords = ($("input[name=word]").val() == '¿Qué buscas?') ? '' : $("input[name=word]").val();
	var categoria = $("select[name=catMenu]").attr("value");
	var estado = $("select[name=estadoMenu]").attr("value");
	var tipo = $("select[name=tipoMenu]").attr("value");     
	pagina ++;
	$.get("http://movil.ebfmex-pub.appspot.com/search",{ pagina:pagina, keywords:keywords, categoria:categoria, estado:estado , tipo:tipo, kind: 'Oferta'},function(data){
		console.log(pagina);
		if(pagina == '1')
		  $(".ofertCont").html('')
		var ofertas = JSON.parse(data);
		if(ofertas.length >= 1){
		  cargaOfertas = true;
		  for(var i in ofertas){
			urlOferta = 'http://home.ebfmex-pub.appspot.com/busqueda-de-ofertas.html';
			addOferta = '<div class="oferta bgWh">'
			addOferta += '<a href="#" class="lighter">'
			addOferta += '<span class="imgcont">'
			addOferta += '<img src="http://home.ebfmex-pub.appspot.com'+ofertas[i].Logo+'" width="212" height="218" alt="'+ofertas[i].Oferta+'" title="'+ofertas[i].Oferta+'" />'
			addOferta += '</span>'
			addOferta += '<h3>'+ofertas[i].Oferta+'</h3>'
			addOferta += '</a>'
			addOferta += '<div class="col-30PR first" style="">'
			addOferta += '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El buen fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer/sharer.php?u=http://localhost/ebfmex/static/busqueda-de-ofertas.html">'
			addOferta += '<img src="../imgs/ofrtTemp/mtShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-40PR first" style="margin-top:5px;">'
			addOferta += '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=' + ofertas[i].Logo + '&p[title]= ' + ofertas[i].Oferta +'">'
			addOferta += '<img src="../imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-30PR first">'
			addOferta += '<a onClick="window.open(\'https://twitter.com/intent/tweet?text=Viendo \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">'
			addOferta += '<img src="../imgs/ofrtTemp/twShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '</div>';
			$(".ofertCont").append(addOferta);                        
			lighterAjax();
			//console.log(i + "_" + j + ': ' + ofertas[i][j] );
		  }    
		}else{
			cargaOfertas = false;
		}
	});
}

