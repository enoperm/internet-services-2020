function makePasswordEnable() {
    var check = document.getElementById("passi").disabled;
    document.getElementById("passi").disabled = !check;
}

function CheckPassword(inputtxt) {
    var passw= /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,}$/;
    return !!inputtxt.match(passw);
}

function getMonthFromString(mon) {
    return new Date(Date.parse(mon +" 1, 2012")).getMonth()+1
 }

function resetStats() {
    var today = new Date().toDateString().split(" ");
    today = today[3] + '-' + getMonthFromString(today[1]) + '-' + today[2]
    console.log(today);

    $('#date-time').value = today;
}

