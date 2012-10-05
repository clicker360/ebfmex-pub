$(document).ready(function() {
	/* Este código comentado debe ir en el template html de lo contrario el 
	 * programa no puede planchar las variables de entorno
	 */
	/*$("#urlimg").attr('href', "{{with .FormDataOf}}{{.Url|js}}{{end}}");*/
	/*var idoft = "{{with .FormDataOf}}{{.IdOft|js}}{{end}}";*/
	$('#loader').hide();
	$('#urlreq').hide();
	$("#url").blur(function() {
		if($('#enlinea').attr('checked') && $('#url').val()=='') { $('#urlreq').show(); } else {$('#urlreq').hide();}
	});
	$("#enviardata").submit(function() {
		if($('#enlinea').attr('checked') && $('#url').val()=='') { $('#urlreq').show(); return false; } else {$('#urlreq').hide(); return true;}
	});
	$("#cards").click(function () { $("#allcards").toggle("slow"); });  
	$("#ONL").click(function () { $("#onlineURL").toggle("slow"); });  
	$("#SUC").click(function () { $("#sucursales").toggle("slow"); });  
	$("#enviar").validationEngine({promptPosition : "topRight", scroll: false});
	$("#enviardata").validationEngine({promptPosition : "topRight", scroll: false});
	var $pic = $("#pic");
	var $urlimg = $("#urimg");
	var max_size=400;
	if(idoft == 'none') {
		// se ocultan los campos que requieren IdOft
		$('#dataoferta').hide();
		$('#imgform').hide();
		$('#modbtn').hide();
		$('#newbtn').show();
	} else {
		$('#dataoferta').show();
		$('#imgform').show();
		$('#modbtn').show();
		$('#newbtn').hide();
	}
	$pic.error(function () { $(this).unbind("error").attr("src", "imgs/imageDefault.jpg"); });
	function updateimg() {
		if(idoft != 'none') {
			var query = "id="+idoft;
			$pic.attr('src', "/ofimg?"+query);
		} else {
			$pic.attr('src', "imgs/imageDefault.jpg");
		}
	}
	updateimg();
	fillpromo();
	(function() {
		
	var bar = $('.bar');
	var percent = $('.percent');
	var status = $('#status');
	   
	$('#enviar').ajaxForm({
		beforeSend: function() {
			status.empty();
			var percentVal = '0%';
			bar.width(percentVal)
			percent.html(percentVal);
		},
		uploadProgress: function(event, position, total, percentComplete) {
			var percentVal = percentComplete + '%';
			bar.width(percentVal)
			percent.html(percentVal);
		},
		complete: function(xhr) {
			status.html(xhr.responseText);
			updateimg();
		}
	}); 

	})();       
	$('textarea[maxlength]').live('keyup blur', function() {
		var maxlength = $(this).attr('maxlength'); var val = $(this).val();
		if (val.length > maxlength) {
			$(this).val(val.slice(0, maxlength));
		}
	});
	$('.promo').click(function(e){
		var token = $(this);
		var stoppick = false;
		if(token.attr("class") == "sugestWord promo") {
			var kids = $('#unpickpromo').children();
			kids.each(function(ii,obj) {
				if($(this).attr('tipo') == token.attr('tipo')) {
					alert("Sólo se permite combinar 3 tipos de promoción diferentes");
					stoppick = true;
				}
			});
			if(!stoppick) {
				$.get("/modpromo", { token: ""+token.attr('value')+"", tipo: ""+token.attr('tipo')+"", id: ""+idoft+"" }, function(resp) {
					if(resp.status=="ok") {
						token.attr("class", "wordselected");	
						$('#unpickpromo').append(token);
					} else {
						alert("Hubo un problema de conexión. Intente de nuevo agregar la promoción");
					}
				}, "json");
			}
		} else {
			$.get("/delpromo", { token: '', tipo: ""+token.attr('tipo')+"", id: ""+idoft+"" }, function(resp) {
				if(resp.status=="ok") {
					token.attr("class", "sugestWord");	
					$('#pickpromo').append(token);
				} else {
					alert("Hubo un problema de conexión. Intente de nuevo eliminar la promoción");
				}
			}, "json");
		}
	});
	$('.pcve').click(function(e){
		var token = $(this);
		if($(this).attr("value") == "0") {
			$.get("/addpcve", { token: ""+token.text()+"", id: ""+idoft+"" }, function(resp) {
				if(resp.status="ok") {
					token.attr("class", "wordselected");	
					token.attr("value", resp.id);	
					$('#unpickpcve').append(token);
				} else {
					alert("Hay problemas de conexión. Intente de nuevo agregar la palabra clave");
				}
			}, "json");
		} else {
			$.get("/delpcve", { id: ""+token.attr('value')+"" }, function(resp) {
				if(resp.status="ok") {
					token.attr("class", "sugestWord");	
					token.attr('value','0');
					$('#pickpcve').append(token);
				} else {
					alert("Hay problemas de conexión. Intente de nuevo eliminar la palabra clave");
				}
			}, "json");
		}
	});
	/* Promociones */
	/* Palabras clave */
	/* Tarjetas participantes */
	
});/* termina onload */

function fillpcve(idoft) {
	var tc=$.get("/pcvesxo", { id: "" + idoft + ""}, function(data) {
		data=$.parseJSON(data);
		alert(data.desc);
		$.each(data, function(i,item){
			alert(item);
		});
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}
function fillpromo() {
	var tc=$.get("/promosxo", { id: "" + idoft + ""}, function(data) {
		var pick;
		var kids = $('#pickpromo').children();
		$.each(data, function(i,item){
			kids.each(function(ii,obj) {
				if($(this).attr('value') == item.token) {
					pick = $(this);
					pick.attr("class", "wordselected");	
					pick.attr("value", item.token);	
					$('#unpickpromo').append(pick);
				}	
			});
		});
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}
function filltc() {
	var tc=$.get("/tarjetasxo", { id: "" + idoft + ""}, function(data) {
		data=$.parseJSON(data);
		alert(data.desc);
		$.each(data, function(i,item){
			alert(item);
		});
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}	
function activateCancel(){ $("#cancelbtn").addClass("show") }
function deactivateCancel(){ $("#cancelbtn").removeClass("show") }

