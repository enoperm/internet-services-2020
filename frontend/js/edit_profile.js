function makePasswordEnable(){
    var check = document.getElementById("passi").disabled;
    if(check)
        document.getElementById("passi").disabled = false;
    else
        document.getElementById("passi").disabled = true;
}

function CheckPassword(inputtxt) 
{ 
    var passw= /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}$/;
    if(inputtxt.match(passw)){ 
      return true;
    }
    else{
      return false;
    }
}

function getMonthFromString(mon){
    return new Date(Date.parse(mon +" 1, 2012")).getMonth()+1
 }

function resetStats(){
    //var username = document.getElementById("unamei").value;
    //var password = document.getElementById("passi").value;
    var today = new Date().toDateString().split(" ");
    today = today[3] + '-' + getMonthFromString(today[1]) + '-' + today[2]
    console.log(today);
    //var dd = String(today.getDate()).padStart(2, '0');
    //var mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
    //var yyyy = today.getFullYear();
    //var lasttime = yyyy + "-" + mm + "-" + dd;

    var formData = {"last_smoke":today,"daily_average":$('#perday').val(),"sticks_per_pack":$('#pack').val(),"price_per_pack":$('#cost').val(),"start_year":$('#year').val()};
    formData = JSON.stringify(formData);

    $.ajax({
        type: 'POST',
        url: '/api/profile',
        data: formData,
        success: function( resp ) {
            console.log( resp );
            console.log('asd');
            $('#date-time').val(today);
        },
        error: function( req, status, err ) {
            console.log( 'something went wrong', status, err );
        }
    });
}

$( '#profile_form' ).submit(function( event ) {
    event.preventDefault();

    var form = $( this );
    var formData = {"last_smoke":$('#date-time').val(),"daily_average":$('#perday').val(),"sticks_per_pack":$('#pack').val(),"price_per_pack":$('#cost').val(),"start_year":$('#year').val()};
    formData = JSON.stringify(formData);

    $.ajax({
        type: 'POST',
        url: '/api/profile',
        data: formData,
        success: function( resp ) {
            console.log( resp );
            console.log("posted");
        }
    });
});

$(function() {
    $.ajax({
        type: 'GET',
        url: '/api/profile',
        success: function( resp ) {
            console.log( resp );
            $('#date-time').val(resp["last_smoke"]);
            $('#perday').val(resp["daily_average"]);
            $('#pack').val(resp["sticks_per_pack"]);
            $('#cost').val(resp["price_per_pack"]);
            $('#year').val(resp["start_year"]);
        },
        error: function( req, status, err ) {
            console.log( 'something went wrong', status, err );
        }
    });
});