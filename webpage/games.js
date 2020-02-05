function fillInOptions(button, options) {
    $.each(options, function (id, val) {
        $('.player-options').append(val);
    });
    $('#game_input_form').attr('game-id', button.data('id'))
    $('#player1').attr('cur', button.data('t1pd'))
    $('#player1').val(button.data('t1pd'))
    $('#player2').attr('cur', button.data('t1po'))
    $('#player2').val(button.data('t1po'))
    $('#player3').attr('cur', button.data('t2pd'))
    $('#player3').val(button.data('t2pd'))
    $('#player4').attr('cur', button.data('t2po'))
    $('#player4').val(button.data('t2po'))
    $('#halfScoreTeam1').val(button.data('t1half'))
    $('#halfScoreTeam1').attr('cur', button.data('t1half'))
    $('#halfScoreTeam2').val(button.data('t2half'))
    $('#halfScoreTeam2').attr('cur', button.data('t2half'))
    $('#endTeam1').val(button.data('t1final'))
    $('#endTeam1').attr('cur', button.data('t1final'))
    $('#endTeam2').val(button.data('t2final'))
    $('#endTeam2').attr('cur', button.data('t2final'))
}

// Fills in the game edit modal with info about the game being viewed. This shows up in a form
// that is editable in order to correct information about the game.
$('#game-edit-modal').on('show.bs.modal', function (event) {
    $.getJSON("players", function (data) {
        var items = [];
        console.log(data)
        $.each(data, function (id, val) {
            items.push("<option id='" + val.ID + "'>" + val.Name + "</option>");
        });
        console.log(items)
        fillInOptions($(event.relatedTarget), items)
    });
});

// Below functions change the class of the Save button in the modal to show that there is something
// to save to the db
$('.edit-field').on("change", function (e) {
    if ($(this).val() != $(this).attr('cur')) {
        $('#save-game-edits').removeClass('btn-outline-primary')
        $('#save-game-edits').addClass('btn-primary')
    } else {
        $('#save-game-edits').removeClass('btn-primary')
        $('#save-game-edits').addClass('btn-outline-primary')
    }
});

// Should only be called from the modal Save button.
function getGameJson() {
    var json = {
        ID: $('#game_input_form').attr('game-id'),
        T1pd: $('#player1').val(),
        T1po: $('#player2').val(),
        T2pd: $('#player3').val(),
        T2po: $('#player4').val(),
        T1half: parseInt($('#halfScoreTeam1').val()),
        T2half: parseInt($('#halfScoreTeam2').val()),
        T1final: parseInt($('#endTeam1').val()),
        T2final: parseInt($('#endTeam2').val())
    };
    return JSON.stringify(json);
};

// Checks if anything changed in the game data and if so sends a PUT request with
// the changed fields to be saved. After the game is successfully changed, the page
// should reload (force a refresh in the browser?)
$('#save-game-edits').on("click", function (e) {
    // check if they changed anything
    if ($(this).hasClass('btn-primary')) {
        var xhttp = new XMLHttpRequest();
        xhttp.onreadystatechange = function () {
            // check if the request is done and returned a OK response
            if (this.readyState == 4) {
                if (this.status == 200) {
                    $('#alert-message').text('Game Saved!');
                    $('#saved-alert').show();
                    $('#game-edit-modal').attr('dirty', 1)
                } else {
                    $('#alert-message').text('Error Saving Game.');
                    $('#saved-alert').removeClass('alert-success').addClass('alert-danger').show();
                }
            }
        };
        xhttp.open("PUT", "edit_game", true)
        // xhttp.setRequestHeader("Content-type", "application/json");
        xhttp.send(getGameJson());
    } else if ($(this).hasClass('btn-success')) {
        $('#game-edit-modal').modal('hide');
    }
});

$('#game-edit-modal').on('hide.bs.modal', function (e) {
    if ($(this).attr('dirty') == 1) {
        location.reload();
    }
    $('#saved-alert').hide();
});

$(document).ready(function() {
    $('#saved-alert').hide();
});