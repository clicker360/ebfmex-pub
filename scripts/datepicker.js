var camp;
var activeDate = 0;
var campoToWrite;
var day1 = 1;
var day2 = 31;
function untogleCalendar(){
	$("#datePick").toggle("slow");
	activeDate = 0;
}
function toggleCalendar(elem){
	campoToWrite = elem;
	campID=$(elem).attr('id');
	if(campID == "date1"){ 
		$("#datePick").css('top', 750);//esta la altura del campo1
	}else{
		$("#datePick").css('top', 800);//esta la altura del campo2
	}
	if(activeDate == 0){
		$("#datePick").toggle("slow");
	}
	activeDate = 1;
	return false;
}
function getTheDate(elem){
	day = $(elem).attr('id');
	if(campID == "date1"){
		day1 = day;
	}else{
		day2 = day;
	}
	if (day1 > day2){
		//$('#alrtDay').show();
	}else{
		$(campoToWrite).val(day+' Nov');
		untogleCalendar();
		//$('#alrtDay').hide();
	}
	return false;
}
