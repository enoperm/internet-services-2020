function check() {
    let message = document.getElementById('message');
    let submitBtn = document.getElementById("next_button");
    
    let pw = document.getElementById('mypassword').value;
    
    if(pw === "") {
        // edge case because message is not styled yet.
        message.textContent = "";
        submitBtn.disabled = true;
        return
    }
    
    if (pw === document.getElementById('mypassword_again').value) {
      let pwRejected = !CheckPassword(pw);
      if(pwRejected) {
        message.textContent = "Password must be at least 8 characters long and contain lowercase, uppercase and numeric characters.";
      } else {
          message.textContent = "";
      }
      submitBtn.disabled = pwRejected;
      return;
    }
    submitBtn.disabled = true;
    message.style.color = 'red';
    message.textContent = 'Passwords do not match';
}

function CheckPassword(inputtxt) {
  let passw =  /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}$/;
  return !!inputtxt.match(passw)
}

$('form').submit(function(e) {
    let passwordValue = document.getElementById("mypassword").value;
    let acceptable = CheckPassword(passwordValue);

    if(!acceptable) {
        document.getElementById("error_message").innerHTML = "";
        e.preventDefault();
        return false;
    }
});