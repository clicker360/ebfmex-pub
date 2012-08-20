function validate(){
		var regmail = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
		var regname = /^[a-zA-Z ]+$/;
		var regtel = /^[0-9\-]+$/;
		var regpass = /^[a-zA-Z \-@#"'()\[\]?¿!¡]+$/;
		var valid = true;
    	//alert(document.formregistro.email1.value);
    	if( !regmail.test(document.formregistro.email.value) ){
			document.formregistro.email.className += " invalid";
			valid = false;
    	}
    	if( !regmail.test(document.formregistro.emailalt.value) ){
			document.formregistro.emailalt.className += " invalid";
			valid = false;
    	}
    	if( !regname.test(document.formregistro.nombre.value) ){
			document.formregistro.nombre.className += " invalid";
			valid = false;
    	}
    	if( !regname.test(document.formregistro.apellidos.value) ){
			document.formregistro.apellidos.className += " invalid";
			valid = false;
    	}
    	if( !regtel.test(document.formregistro.tel.value) ){
			document.formregistro.tel.className += " invalid";
			valid = false;
    	}
    	if( !regtel.test(document.formregistro.cel.value) ){
			document.formregistro.cel.className += " invalid";
			valid = false;
    	}
    	if( !regpass.test(document.formregistro.pass.value) ){
			document.formregistro.pass.className += " invalid";
			valid = false;
    	}
    	if( !regpass.test(document.formregistro.pass2.value) ){
			document.formregistro.pass2.className += " invalid";
			valid = false;
    	} else {
    		if( document.formregistro.pass.value != document.formregistro.pass2.value ){
				document.formregistro.pass2.className += " invalid";
				valid = false;   	
    		}
    	}
    	if ( !document.formregistro.terminos.checked ){
			document.formregistro.terminos.className += " invalid";
    		valid = false;  
    	}
    	if ( document.formregistro.puesto.value == "" ){
			document.formregistro.puesto.className += " invalid";
    		valid = false;  
    	}
    	return valid;
	}