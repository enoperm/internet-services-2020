$( 'form' ).submit(function( event ) {
    event.preventDefault();

    var form = $( this );
    var formData = {"username":$('#username').val(),"password":$('#password').val()};
    formData = JSON.stringify(formData);
    //console.log(formData);

    $.ajax({
        type: 'POST',
        url: '/api/session',
        data: formData,
        success: function( resp ) {
            console.log( resp );
            window.location.href = '/main.html';
        },
        error: function( req, status, err ) {
            console.log( 'something went wrong', status, err );
            if (req.status == 401) {
                $('#errormessage').text('Failed, wrong credentials.');
            }
        }
    });
});