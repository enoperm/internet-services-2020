function check() {
    if (document.getElementById('mypassword').value == document.getElementById('mypassword_again').value) {
      document.getElementById('message').style.color = 'green';
      document.getElementById('message').innerHTML = 'Password is matching';
      document.getElementById("next_button").disabled = false;
    } else {
      document.getElementById('message').style.color = 'red';
      document.getElementById('message').innerHTML = 'Password is not matching';
      document.getElementById("next_button").disabled = true;
    }
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

$( 'form' ).submit(function( event ) {
    event.preventDefault();

    var form = $( this );
    var formData = {"username":$('#uniqueusername').val(),"password":$('#mypassword').val()};
    formData = JSON.stringify(formData);
    console.log(formData);

    var passwordValue = document.getElementById("mypassword").value;
    var check = CheckPassword(passwordValue);

    if(check){
        $.ajax({
            type: 'POST',
            url: '/api/register',
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