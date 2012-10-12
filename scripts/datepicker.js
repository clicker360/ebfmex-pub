var camp;
var campoToWrite;
function untogleCalendar(){
	$("#datePick").toggle("slow");
	activeDate = 0;
}
function toggleCalendar(elem){
	campoToWrite = elem;
	campID=$(elem).attr('id');
	$("#datePick").css('top', 550);//esta la altura del campo1
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
$("#sugerencia").mouseover(function(event){
	var elmn = $("#slug1");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
  	elmn.css('left', theleft+40);
	elmn.css('top', thetop-100);
	elmn.toggle("fast");    
   });
   
$("#sugerencia").mouseout(function() { 
   $("#lsug1").toggle("fast");
   });

$("#dater").mouseover(function(event){
	var elmn = $("#slug2");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
	elmn.css('left', theleft+20);
	elmn.css('top', thetop-30);
	elmn.toggle("fast");
   });
   
$("#dater").mouseout(function() { 
 $("#slug2").toggle("fast");
   });
   
$("#worder").mouseover(function(event){
	var elmn = $("#slug3");
	var posicion = $(this).position();
	var theleft= posicion.left + $(this).width();
	var thetop =posicion.top;
	elmn.css('left', theleft+20);
	elmn.css('top', thetop-30);
	elmn.toggle("fast");
   });
   
$("#worder").mouseout(function() { 
 $("#slug3").toggle("fast");
   });
});

