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

function resetStats(){
    var username = document.getElementById("unamei").value;
    var password = document.getElementById("passi").value;

    var today = new Date();
    var dd = String(today.getDate()).padStart(2, '0');
    var mm = String(today.getMonth() + 1).padStart(2, '0'); //January is 0!
    var yyyy = today.getFullYear();

    var lasttime = yyyy + "-" + mm + "-" + dd;

}

$( 'form' ).submit(function( event ) {
    event.preventDefault();

    var form = $( this );
    var formData = {"username":$('#unamei').val(),"password":$('#passi').val(), "perday":$('#perday').val(),"datetime":$('#date-time').val(),
                    "perday":$('#perday').val(),"pack":$('#pack').val(), "cost":$('#cost').val(), "year":$('#year').val()};
    formData = JSON.stringify(formData);

    var passwordValue = document.getElementById("passi").value;
    var check = CheckPassword(passwordValue);

    if(check){
        $.ajax({
            type: 'POST',
            url: '/api/editprofile',
            data: formData,
            success: function( resp ) {
                console.log( resp );
                window.location.href = '/';
            }
        });
    }
    else{
        document.getElementById("error_message").innerHTML="Wrong password format";
    }
});