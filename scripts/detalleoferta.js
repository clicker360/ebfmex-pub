$(document).ready(function(){
    var queryVars = getVars();
    id = (queryVars.hasOwnProperty('id')) ? queryVars.id : false;
    hasA = false
    posA = false
    for(var i in id){
        if(id[i] == '#' && !hasA){
            hasA = true;
            posA = i;
        }
    }
    id = (posA) ? id.substring(0,posA) : id;
    if(!id)
        window.location.href = '/';
    getOferta(id);
    // Creating a new map

    $("#buscar").submit(function(){
        if($("input[name=word]").val() == '¿Qué buscas?'){
            $("input[name=word]").val('')
        }
    })

	
});
function sucMap(lat, lng, div){
    var zoom = 17;
    var marker = false;
	if(lat == '0') {
		lat = 22.770856;
		lng = -102.583243;
		zoom = 4;
	}
	var center = new google.maps.LatLng(lat,lng);
	var options = {
		zoom: zoom,
		center: center,
		mapTypeId: google.maps.MapTypeId.ROADMAP,
		streetViewControl: false
	};
	map = new google.maps.Map(document.getElementById(div), options);
	if (!marker) {
		// Creating a new marker and adding it to the map
		marker = new google.maps.Marker({
			map: map,
			draggable: true
		});
		marker.setPosition(center);
	}
	google.maps.event.addListener(marker, 'dragend', function() {
		var tmppos = ''+this.getPosition();
		var latlng = tmppos.split(',');
		lat = latlng[0].replace('(','');
		lng = latlng[1].replace(')','')
		map.setCenter(this.getPosition());
	});
        // Creating a new map    
}
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
function getOferta(id){
    $.get('http://pruebas.ebfmxorg.appspot.com/wsdetalle',{
        id:id
    },function(oferta){
        if(typeof(oferta) != 'object')
            oferta = JSON.parse(oferta);
        console.log(oferta);
        var urlOferta = 'http://www.elbuenfin./detalleoferta.html?id='+oferta.idoft;
        $("#logoOft img").attr('src','http://pruebas.ebfmxorg.appspot.com/simg?id='+oferta.idemp);
        $("#nomEmp h4").html(oferta.empresa);
        $("#titOft h3").html(oferta.oferta);
        $("#desOft p").html(oferta.descripcion);
        $("#msEmp").attr('href','/micrositio.html?id='+oferta.idemp);
        $("#imgOft img").attr('src','http://pruebas.ebfmxorg.appspot.com/ofimg?id='+oferta.imgurl);
        $("#mtShare").html('<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El Buen Fin \' + this.href, this.target, \'width=600,height=400\'); return false;" href="'+urlOferta+'"><img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" /></a>')
        $("#fbShare").html('<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=http://pruebas.ebfmxorg.appspot.com/ofimg?id='+oferta.imgurl+'&p[title]= ' + oferta.oferta +'&p[summary]='+oferta.descripcion+'"><img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" /></a>')
        $("#twShare").html('<a onClick="window.open(\'https://twitter.com/intent/tweet?text='+oferta.empresa+ ' '+ oferta.oferta+' \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter"><img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" /></a>');

		$("img").error(function() {
			console.log("img=none");
			img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='212' height='218'/>";
			$(this).replaceWith(img);
		});


        if(oferta.url)
            $("#enLinea").html('<div class="col-12 bgRd marg-B10px marg-T70px padd-R10px marg-L5px"><h4 class=" typ-Wh"> El Buen Fin en Línea</h4></div><div><a target="_blank" href="'+oferta.url+'" id="urlOft" >'+oferta.url+'</a></div>');
        getSucursales(oferta.idemp);
    })
}
function getSucursales(id){
    $.get('http://pruebas.ebfmxorg.appspot.com/wssucs',{
        id:id
    },function(sucursales){
         if(typeof(sucursales) != 'object')
            sucursales = JSON.parse(sucursales);
        for(var i in sucursales){
            $("#sucList").append('<li><a href="#null" onClick="showMap(\''+sucursales[i].lat+'\',\''+sucursales[i].lng+'\',\'map\'); return false;">'+sucursales[i].sucursal+'</a></li>');
            console.log(sucursales[i]);
        }
    })
}
function showMap(lat, lng, div){
    //if(lat != '0' && lng != '0'){
        if($("#mapCont").is(":visible"))
            $('#mapCont').slideToggle('slow', function() {});        
        sucMap(lat , lng , div);
        $('#mapCont').slideToggle('slow', function() {});
    //}
}
