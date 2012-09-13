// Defining some global variables
var map, geocoder, marker, infowindow;
$(document).ready(function() {
	// Creating a new map
	var zoom = 17;
	var lat = $('#lat').val();
	var lng = $('#lng').val();
	if(lat == '') { 
		lat = 22.770856;
		lng = -102.583243; 
		zoom = 4;
	}
	var center = new google.maps.LatLng(lat,lng);
	var options = {
		zoom: zoom,
		center: center,
		mapTypeId: google.maps.MapTypeId.ROADMAP
	};
	map = new google.maps.Map(document.getElementById('map'), options);
	if (!marker) {
		// Creating a new marker and adding it to the map
		marker = new google.maps.Marker({
			map: map,
			draggable: true
		});
		marker.setPosition(center);
	}
	// Getting a reference to the HTML form
	$('#buscar').bind('keydown keyup mousedown',function(e){
		locateAddress();
	});
	google.maps.event.addListener(marker, 'dragend', function() {
		var tmppos = ''+this.getPosition();
		var latlng = tmppos.split(',');
		document.getElementById('lat').value = latlng[0].replace('(','');
		document.getElementById('lng').value = latlng[1].replace(')','')
		map.setCenter(this.getPosition());
	});
});

function locateAddress() {
	// Getting the address from the text input
	var estado = document.getElementById('estado').value;
	var municipio = document.getElementById('municipio').value;
	var calle = document.getElementById('calle').value;
	var colonia = document.getElementById('colonia').value;
	var cp = document.getElementById('cp').value;
	var address = calle + ", " + colonia + ", " + cp + ", " + municipio + ", " + estado + ", MEXICO";
	
	// Check to see if we already have a geocoded object. If not we create one
	if(!geocoder) {
		geocoder = new google.maps.Geocoder();
	}
	// Creating a GeocoderRequest object
	var geocoderRequest = {
		address: address
	}
	// Making the Geocode request
	geocoder.geocode(geocoderRequest, function(results, status) {
		// Check if status is OK before proceeding
		if (status == google.maps.GeocoderStatus.OK) {
			// Center the map on the returned location
			map.setCenter(results[0].geometry.location);
			map.setZoom(17);
			// Check to see if we've already got a Marker object
			if (!marker) {
				// Creating a new marker and adding it to the map
				marker = new google.maps.Marker({
					map: map,
					draggable: true
				});
			}
			// Setting the position of the marker to the returned location
			marker.setPosition(results[0].geometry.location);
			// Check to see if we've already got an InfoWindow object
			/*if (!infowindow) {
				// Creating a new InfoWindow
				infowindow = new google.maps.InfoWindow();
			}
			// Creating the content of the InfoWindow to the address
			// and the returned position
			var content = '<strong>' + results[0].formatted_address + '</strong><br />';
			content += 'Lat: ' + results[0].geometry.location.lat() + '<br />';
			content += 'Lng: ' + results[0].geometry.location.lng();
			// Adding the content to the InfoWindow
			infowindow.setContent(content);
			// Opening the InfoWindow
			infowindow.open(map, marker);*/
			document.getElementById('lat').value = results[0].geometry.location.lat();
			document.getElementById('lng').value = results[0].geometry.location.lng();
			
			
		}
	});
}

