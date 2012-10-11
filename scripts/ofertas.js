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
	$("#enviar").validationEngine({promptPosition : "topRight", scroll: false});
	$("#enviardata").validationEngine({promptPosition : "topRight", scroll: false});
	var $pic = $("#pic");
	var $urlimg = $("#urimg");
	var max_size=400;
	if(idoft == 'none') {
		// se ocultan los campos que requieren IdOft
		putDefault();
		$('#dataoferta').hide();
		$('#imgform').hide();
		$('#modbtn').hide();
		$('#newbtn').show();
		$('#statuspub').attr("checked", true);
	} else {
		$('#dataoferta').show();
		$('#imgform').show();
		$('#modbtn').show();
		$('#newbtn').hide();	
	}
	
	
	/*$pic.error(function () { $(this).unbind("error").attr("src", "imgs/imageDefault.jpg"); });*/
	

	/* solo se actualizan estos datos si hay id de oferta */
	if(idoft != 'none') {
		updateimg(idoft);
		fillpromo();
		//filltc(idoft);
		fillpcve(idoft, idemp);
		fillsucursales(idoft, idemp);
	}


	var bar = $('.bar');
	var percent = $('.percent');
	var status = $('#status');
	var img;
		
	   
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
			setTimeout(function(){
				 updateimg(idoft); }, 1000); 	
		}
	}); 
$("#pic").error(function() { putDefault()});
	     

	$('textarea[maxlength]').live('keyup blur', function() {
		var maxlength = $(this).attr('maxlength'); var val = $(this).val();
		if (val.length > maxlength) {
			$(this).val(val.slice(0, maxlength));
		}
	});

	/* Promociones
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
						alert("Hubo un problema de conexión. Intente agregar de nuevo la promoción");
					}
				}, "json");
			}
		} else {
			$.get("/delpromo", { token: '', tipo: ""+token.attr('tipo')+"", id: ""+idoft+"" }, function(resp) {
				if(resp.status=="ok") {
					token.attr("class", "sugestWord");	
					$('#pickpromo').append(token);
				} else {
					alert("Hubo un problema de conexión. Intente agregar de nuevo la promoción");
				}
			}, "json");
		}
	});
 */
	/* Palabras clave */
	$('#pcvepicker').on("click", "a", function(e){
		var token = $(this);
		if($(this).attr("value") == "0") {
			$.get("/addword", { token: ""+token.text()+"", id: ""+idoft+"" }, function(resp) {
				if(resp.status=="ok") {
					token.attr("class", "wordselected");	
					token.attr("value", resp.id);	
					$('#unpickpcve').append(token);
				} else {
					alert("Hay problemas de conexión. Intente agregar de nuevo la palabra clave");
				}
			}, "json");
		} else {
			$.get("/delword", { id: ""+token.attr('value')+"", token: ""+token.text()+"" }, function(resp) {
				if(resp.status=="ok") {
					token.attr("class", "sugestWord");	
					token.attr('value','0');
					$('#pickpcve').append(token);
				} else {
					alert("Hay problemas de conexión. Intente eliminar de nuevo la palabra clave");
				}
			}, "json");
		}
	});
	$("#nuevapcve").click(function(e) {
		var token = $("#tokenpcve");
		$.get("/addword", { token: ""+token.val()+"", id: ""+idoft+"" }, function(resp) {
			if(resp.status=="ok") {
				clearpcve();
				fillpcve(idoft,idemp);
			} else {
				alert("Hay problemas de conexión. Intente agregar de nuevo la palabra clave");
			}
		}, "json");
	});

});/* termina onload */

function avoidCache(){
				var numRam = Math.floor(Math.random() * 500);
				return numRam;
			}		
	function putDefault() {
				$('#pic').remove();
				img = "<img  src = 'imgs/imageDefault.jpg' id='pic' width='258px' />" 
				$('#urlimg').append(img);
			}
	function updateimg(idoft) {
		if(idoft != 'none') {
			//alert(idoft);
			$('#pic').remove();
			var query = "id="+idoft + "&Avc=" + avoidCache();
			img = "<img  src = '/ofimg?"+ query +"' id='pic' width='256px' />" 
			$('#urlimg').append(img);
		} else {
	 		putDefault();
		}
	}


/*
 * Llena palabras clave por oferta y empresa
 */
function fillpcve(idoft, idemp) {
	$.get("/wordsxo", { id: "" + idoft + ""}, function(data) {
		if($.isArray(data)) {
			$.each(data, function(i,item){
				var anchor = "<a href=\"#null\" class=\"wordselected\" id=\"pcve_"+item.token+"\" value=\""+item.id+"\">"+item.token+"</a>"
				$('#unpickpcve').append(anchor);
			});
		}
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});

	$.get("/wordsxe", { id: "" + idemp + ""}, function(data) {
		if($.isArray(data)) {
			$.each(data, function(i,item){
				// Si en el ajax anterior no se añadio algo, aquí se añade como no seleccionado
				if($("#pcve_"+item.token).length == 0) {
					var anchor = "<a href=\"#null\" class=\"sugestWord pcve\" id=\"pcve_"+item.token+"\" value=\"0\">"+item.token+"</a>"
					$('#pickpcve').append(anchor);
				}
			});
		}
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});

}
function clearpcve() {
	$("#unpickpcve").empty();
	$("#pickpcve").empty();
}

/*
 * Llena las promociones
 */
function fillpromo() {
	var tc=$.get("/promosxo", { id: "" + idoft + ""}, function(data) {
		var pick;
		var kids = $('#pickpromo').children();
		if($.isArray(data)) {
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
		}
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}

/*
 * Llena tarjetas participantes
 
function filltc(idoft) {
	$.get("/tcxo", { id: "" + idoft }, function(data) {
		$.each(data, function(i,item){
			var div = "<div class=\"bg-Gry2 col-10\"><label class=\"col-4 marg-L5pix\">"+item.Tarjeta+"</label> <input name=\"tarjeta\" type=\"checkbox\" class=\"last marg-R15pix marg-U5pix\" value=\""+item.Tarjeta+"\" id=\"tc_"+item.Id+"\" /></div>";
			$('#listatc').append(div);
			if(item.Selected=="1") {
				$("#tc_"+item.Id).attr('checked', true);
			} else {
				$("#tc_"+item.Id).attr('checked', false);
			}
			$("#tc_"+item.Id).change(function() { 
				var tarjetas = $("#listatc").children();
				var sep = "";
				var jsontxt = "";
				tarjetas.each(function(ii,obj) {
					i = ii+1;
					if($("#tc_"+i).is(':checked')) {
						jsontxt += sep+"{ \"Id\":"+i+", \"Tarjeta\":\""+$("#tc_"+i).val()+"\", \"Selected\":1, \"Status\":\"\" }";
					} else {
						jsontxt += sep+"{ \"Id\":"+i+", \"Tarjeta\":\""+$("#tc_"+i).val()+"\", \"Selected\":0, \"Status\":\"\" }";
					}
					sep = ",";
				});
				$.get("/modtcxo", { token: "["+jsontxt+"]", id: ""+idoft+"" }, function(resp) {
					if(resp.Status=="ok") {
						//alert(resp.status);
					} else {
						alert("Hay problemas de conexión. Intente agregar de nuevo la palabra clave");
					}
				}, "json");
			});
		});
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}
 * Llena las sucursales
 */
function fillsucursales(idoft, idemp) {
	$.get("/ofsuc", { idoft: "" + idoft + "", idemp: "" + idemp + ""}, function(data) {
		$.each(data, function(i,item){
			var div = "<div class=\"gridsubRow bg-Gry2\"><label class=\"col-5 marg-L10pix\">"+item.sucursal+"</label><input name=\""+item.idsuc+"\" type=\"checkbox\" class=\"last marg-5px\" id=\""+item.idsuc+"\"/></div>";
			$('#listasuc').append(div);
			if(item.idoft!="") {
				$("#"+item.idsuc).attr('checked', true);
			} else {
				$("#"+item.idsuc).attr('checked', false);
			}
			$("#"+item.idsuc).change(function() { 
				if($(this).is(':checked')) {
					$.get("/addofsuc", { idoft: "" + idoft + "", idemp: "" + idemp + "", idsuc: "" + item.idsuc + ""}, function(data) { })
					.success(function(){})
					.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
				} else {
					$.get("/delofsuc", { idoft: "" + idoft + "", idsuc: "" + item.idsuc + ""}, function(data) { })
					.success(function(){})
					.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
				}
			});
		});
	})
	.success(function(){})
	.error(function(){alert('Hay problemas de conexión, espere un momento y refresque la página');})
	.complete(function(){});
}
function activateCancel(){ $("#cancelbtn").addClass("show") }
function deactivateCancel(){ $("#cancelbtn").removeClass("show") }

