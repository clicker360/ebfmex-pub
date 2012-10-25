$(document).ready(function(){
    $("#CloserLight").click(function() {
		$('#cuerpo').removeClass('noscroll');
		$('#lightback').addClass("hide");
		$('#lightfront').addClass("hide");
        return false;
	});    
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
    getEmpresa(id);
    getSucursales(id);
    pagina = 0;
    getOfertas(id);
    cargaOfertas = false;
    inSearch = false;
    $(document).scroll(function(){
        if(cargaOfertas && !inSearch){
            if(($(document).scrollTop() + $(window).height()) >= ($(document).height())){
                getOfertas(id);
            }
        }
    });
});
function lighterAjax(){
    $(".lighter").click(function() {
        var id = $(this).parent().attr('id');
        $.get('http://pruebas.ebfmxorg.appspot.com/wsdetalle',{
            id:id
        },function(data){
            var imgurl = (data.hasOwnProperty('imgurl')) ? 'http://pruebas.ebfmxorg.appspot.com/ofimg?id='+data.imgurl : '';
            var titOft = (data.hasOwnProperty('oferta')) ? data.oferta : '';
            var desOft = (data.hasOwnProperty('descripcion')) ? data.descripcion : '';
            var nomEmp = (data.hasOwnProperty('empresa')) ? data.empresa : '';
            var idEmp = (data.hasOwnProperty('idemp')) ? data.idemp : '';
            var enLinea = (data.hasOwnProperty('enlinea')) ? (data.enlinea) ? data.url : false : false;
            if(enLinea){
                $("#enLinea").html('<div class="col-12 bgRd marg-B10px marg-T10px padd-R10px marg-L5px" ><h4 class="typ-Wh"> El Buen Fin en Línea</h4></div><div class="col-13 marg-B10px marg-T10px padd-R10px marg-L10px"><a target="_blank" href="'+enLinea+'" class="first" >'+enLinea+'</a></div>')
            }
            var urlOferta = 'http://localhost:8080/detalleoferta.html?id='+data.idoft;
            var mtOft = '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El Buen Fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="'+urlOferta+'">'
            mtOft += '<img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" />'
            mtOft += '</a>'
            var fbOft = '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=' + imgurl + '&p[title]= ' + titOft +'">';
            fbOft += '<img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />';
            fbOft += '</a>';
            var twOft = '<a onClick="window.open(\'https://twitter.com/intent/tweet?text=Viendo \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">';
            twOft += '<img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" />';
            twOft += '</a>';
            $("#sucList").html('');
            if(idEmp){
                var imgEmp = 'http://pruebas.ebfmxorg.appspot.com/simg?id='+idEmp;
                $(".logoOferta img").attr('src',imgEmp);
                $.get('http://pruebas.ebfmxorg.appspot.com/wssucs',{
                    id:idEmp
                },function(sucursales){
                    if(sucursales.length >= 1){
                        for(var i in sucursales){
                            $("#sucList").append('<li><a href="#null" onClick="showMap(\''+sucursales[i].lat+'\',\''+sucursales[i].lng+'\',\'map\',\'mapCont\'); return false;">'+sucursales[i].sucursal+'</a></li>');
                        }
                    }
                    //showMap();
                });
            }
            //$("#divCloseMap").html('<span class="col-9 Dsblock alineCenter marg-L20px" id="closeMap">CERRAR MAPA [X]</span>')
            $("#msEmp a").attr('href','/micrositio.html?id='+idEmp)
            $("#imgOft img.img").attr('src', imgurl);
            $("#titOft h3").html(titOft);
            $("#desOft p").html(desOft);
            $("#nomEmp h4").html(nomEmp);
            $("#mtOft").html(mtOft);
            $("#fbOft").html(fbOft);
            $("#twOft").html(twOft);
            $('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
            $('#lightback').removeClass("hide");
            $('#lightfront').removeClass("hide");
        });
        return false;
    });
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
function getEmpresa(id){
    $.get('http://pruebas.ebfmxorg.appspot.com/wsmicrositio',{id:id},function(empresa){
        var urlEmpresa = 'http://localhost:8080/micrositio.html?id='+empresa.idemp;
        $("#nomEmp h4").html(empresa.name);
        $("#desEmp p").html(empresa.desc);
        var imgEmp = 'http://pruebas.ebfmxorg.appspot.com/simg?id='+empresa.idemp;
        $(".logoOferta img").attr('src',imgEmp);
        var urlEmp = (empresa.url) ? empresa.url : '#';
        $("#urlEmp").attr('href',urlEmp);
        (urlEmp != '#') ? $("#urlEmp").attr('target','_blank') : $("#urlEmp").removeAttr('target');
        $("#mtShareE").html('<a onClick="window.open(\'mailto:?subject=Conoce esta empresa&body=Conoce esta empresa de El Buen Fin \' + this.href, this.target, \'width=600,height=400\'); return false;" href="'+urlEmpresa+'"><img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" /></a>')
        $("#fbShareE").html('<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlEmpresa + '&p[images][0]='+imgEmp+'&p[title]= ' + empresa.name +'"><img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" /></a>')
        $("#twShareE").html('<a onClick="window.open(\'https://twitter.com/intent/tweet?text=Viendo \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlEmpresa +'" class="btwitter" title="Compartelo en Twitter"><img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" /></a>');
    })
}
function getSucursales(id){
    $("#sucListM").html('');
    $.get('http://pruebas.ebfmxorg.appspot.com/wssucs',{id:id},function(sucursales){
        if(sucursales.length >= 1){
            for(var i in sucursales){
                $("#sucListM").append('<li><a href="#null" onClick="showMap(\''+sucursales[i].lat+'\',\''+sucursales[i].lng+'\',\'mapM\',\'mapContM\'); return false;">'+sucursales[i].sucursal+'</a></li>')
            }
        }
        //showMapM();
    })
}
function getOfertas(id){
    pagina++;
    if(pagina == '1')
        $("#ofertCont").html('')
    $("#ofertCont").append('<div class="cargando"><h4 style="float:left; width:100%; text-align:center;" >Cargando...</h4></div>');
    inSearch = true;
    $.get("http://pruebas.ebfmxorg.appspot.com/wsofxe",{id:id,pagina:pagina},function(ofertas){
        console.log(pagina);
        if(ofertas.length >= 1){
            cargaOfertas = true;
            for(var i in ofertas){
                var urlOferta = 'http://localhost:8080/detalleoferta.html?id='+ofertas[i].idoft;
                addOferta = '<div class="oferta bgWh" id="'+ofertas[i].idoft+'">'
                addOferta += '<a href="/detalleoferta.html?id='+ofertas[i].idoft+'" class="lighter">'
                addOferta += '<span class="imgcont">'
                addOferta += '<img src="http://pruebas.ebfmxorg.appspot.com/ofimg?id='+ofertas[i].imgurl+'" width="212" height="218" alt="'+ofertas[i].oferta+'" title="'+ofertas[i].oferta+'" />'
                addOferta += '</span>'
                addOferta += '<h3>'+ofertas[i].oferta+'</h3>'
                addOferta += '</a>'
                addOferta += '<div class="col-30PR first" style="">'
                addOferta += '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El Buen Fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="'+urlOferta+'">'
                addOferta += '<img src="/imgs/ofrtTemp/mtShare.jpg" alt="Enviar por correo electrónico" />'
                addOferta += '</a>'
                addOferta += '</div>'
                addOferta += '<div class="col-40PR first" style="margin-top:5px;">'
                addOferta += '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=http://pruebas.ebfmxorg.appspot.com/ofimg?id='+ofertas[i].imgurl+'&p[title]= ' + ofertas[i].oferta +'">'
                addOferta += '<img src="/imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />'
                addOferta += '</a>'
                addOferta += '</div>'
                addOferta += '<div class="col-30PR first">'
                addOferta += '<a onClick="window.open(\'https://twitter.com/intent/tweet?text=Viendo \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">'
                addOferta += '<img src="/imgs/ofrtTemp/twShare.jpg" alt="Compartir en Twitter" />'
                addOferta += '</a>'
                addOferta += '</div>'
                addOferta += '</div>';
                $("#ofertCont").append(addOferta);
            }
            lighterAjax();
        }else{
            cargaOfertas = false;
            $("#ofertCont").append('<h4 style="float:left; width:100%; text-align:center;">No hay más ofertas de esta empresa.</h4>');
        }
        $(".cargando").remove();
        inSearch = false;
    })
}
function showMap(lat, lng, div, divCont){
    //if(lat != '0' && lng != '0'){
        if($("#"+divCont).is(":visible"))
            $("#"+divCont).slideToggle('slow', function() {});
        sucMap(lat , lng , div);
        $("#"+divCont).slideToggle('slow', function() {});
    //}
}
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
