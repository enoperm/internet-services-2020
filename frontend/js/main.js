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