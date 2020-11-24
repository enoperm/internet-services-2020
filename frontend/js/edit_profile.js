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