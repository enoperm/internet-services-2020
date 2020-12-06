$( '#logout' ).click(function( event ) {
    event.preventDefault();

    var formData = {"username":$('#username').val(),"password":$('#password').val()};
    formData = JSON.stringify(formData);
    console.log(formData);

    $.ajax({
        type: 'DELETE',
        url: '/api/session',
        data: formData,
        success: function( resp ) {
            console.log( resp );
            window.location.href = '/main.html';
        }
    });
});