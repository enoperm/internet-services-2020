var check = function() {
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

function write_to_console() {
    var usernameValue = document.getElementById("uniqueusername").value;
    var passwordValue = document.getElementById("mypassword").value;
    var password_againValue = document.getElementById("mypassword_again").value;
  
    console.log("Username: ", usernameValue, "\n" ,  "Password: ", passwordValue, "\n", "Password again: ", password_againValue);   
}
