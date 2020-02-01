$('#game-edit-modal').on('show.bs.modal', function (event) {
    var button = $(event.relatedTarget) // Button that triggered the modal
    // place holder for ajax query to get player options? Could load with the original page laod instead
    $('#player1 option[value="cur"]').text(button.data('t1pd'))
    $('#player2 option[value="cur"]').text(button.data('t1po'))
    $('#player3 option[value="cur"]').text(button.data('t2pd'))
    $('#player4 option[value="cur"]').text(button.data('t2po'))
    $('#halfScoreTeam1').val(button.data('t1half'))
    $('#halfScoreTeam1').attr('cur', button.data('t1half'))
    $('#halfScoreTeam2').val(button.data('t2half'))
    $('#halfScoreTeam2').attr('cur', button.data('t2half'))
    $('#endTeam1').val(button.data('t1final'))
    $('#endTeam1').attr('cur', button.data('t1final'))
    $('#endTeam2').val(button.data('t2final'))
    $('#endTeam2').attr('cur', button.data('t2final'))
})

$('#halfScoreTeam1').on("change", function(e) {
    if ($(this).val() != $(this).attr('cur')) {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#halfScoreTeam2').on("change", function(e) {
    if ($(this).val() != $(this).attr('cur')) {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#endTeam1').on("change", function(e) {
    if ($(this).val() != $(this).attr('cur')) {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#endTeam2').on("change", function(e) {
    if ($(this).val() != $(this).attr('cur')) {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});

$('#player1').on("change", function(e) {
    if ($(this).val() != "cur") {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#player2').on("change", function(e) {
    if ($(this).val() != "cur") {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#player3').on("change", function(e) {
    if ($(this).val() != "cur") {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});
$('#player4').on("change", function(e) {
    if ($(this).val() != "cur") {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').addClass('btn-outline-primary')
        $('#save-game-edits').removeClass('btn-primary')
    }
});

$('#save-game-edits').on("click", function(e) {
    // placeholder for ajax call to game input route with PUT method
});