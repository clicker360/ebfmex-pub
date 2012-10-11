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
