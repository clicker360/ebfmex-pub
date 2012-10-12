var camp;
var campoToWrite;
function untogleCalendar(){
	$("#datePick").toggle("slow");
	activeDate = 0;
}
function toggleCalendar(elem){
	campoToWrite = elem;
	campID=$(elem).attr('id');
	$("#datePick").css('top', 450);//esta la altura del campo1
	$("#datePick").toggle("slow");
	return false;
}
function getTheDate(elem){
	var day = $(elem).attr('id');
		$(campoToWrite).val(day+' Nov');
		untogleCalendar();
	return false;
}
$(document).ready(function() { 

$("#sugerencia").mouseenter(slug1In).mouseleave(slug1Out);

function slug1In(){
	var elmn = $("#slug1");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
  	elmn.css('left', theleft+40);
	elmn.css('top', thetop-100);
	elmn.toggle(350);    
   }
   
function slug1Out() { 
   $("#slug1").toggle(350);s
   }
   
  $("#dater").mouseenter(slug2In).mouseleave(slug2Out); 
function slug2In(){
	var elmn = $("#slug2");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
	elmn.css('left', theleft+40);
	elmn.css('top', thetop-30);
	elmn.toggle("fast");
   }
   
function slug2Out(){
 $("#slug2").toggle(400);
   }
   
$("#worder").mouseenter( slug3In).mouseleave(slug3Out);  
function slug3In(){
	var elmn = $("#slug3");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
	elmn.css('left', theleft+40);
	elmn.css('top', thetop-20);
	elmn.toggle(400);
   }
   
function slug3Out(){
 $("#slug3").toggle(400);
   }
});

