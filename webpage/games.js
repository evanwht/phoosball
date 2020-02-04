function fillInOptions(button, options) {
    $.each(options, function(id, val) {
        $('.player-options').append(val);
    });
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
    $.getJSON( "players", function( data ) {
        var items = [];
        console.log(data)
        $.each( data, function( id, val ) {
          items.push( "<option id='" + val.ID + "'>" + val.Name + "</option>" );
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
// $('#halfScoreTeam1').on("change", function (e) {
//     if ($(this).val() != $(this).attr('cur')) {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#halfScoreTeam2').on("change", function (e) {
//     if ($(this).val() != $(this).attr('cur')) {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#endTeam1').on("change", function (e) {
//     if ($(this).val() != $(this).attr('cur')) {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#endTeam2').on("change", function (e) {
//     if ($(this).val() != $(this).attr('cur')) {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#player1').on("change", function (e) {
//     if ($(this).val() != "cur") {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#player2').on("change", function (e) {
//     if ($(this).val() != "cur") {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#player3').on("change", function (e) {
//     if ($(this).val() != "cur") {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });
// $('#player4').on("change", function (e) {
//     if ($(this).val() != "cur") {
//         $('#save-game-edits').removeClass('btn-outline-primary')
//         $('#save-game-edits').addClass('btn-primary')
//     } else {
//         $('#save-game-edits').addClass('btn-outline-primary')
//         $('#save-game-edits').removeClass('btn-primary')
//     }
// });

// Should only be called from the modal Save button. Gathers the changed fields 
// in the modal form into a json object to be sent to the data base
function getChangedGameJson() {
    // set of fields to pull from
    // ID, t1pd, t1po, t2pd, t2po, t1half, t2half, t1final, t2final
    var json = {
        id: 1
    };
    if ($('#player1').val() != $('#player1').attr('cur')) {
        json['T1pd'] = $('#player1').val();
    }
    if ($('#player2').val() != $('#player2').attr('cur')) {
        json['T1po'] = $('#player2').val();
    }
    if ($('#player3').val() != $('#player3').attr('cur')) {
        json['T2pd'] = $('#player3').val();
    }
    if ($('#player4').val() != $('#player4').attr('cur')) {
        json['T2po'] = $('#player4').val();
    }
    if ($('#halfScoreTeam1').val() != $('#halfScoreTeam1').attr('cur')) {
        json['T1half'] = parseInt($('#halfScoreTeam1').val());
    }
    if ($('#halfScoreTeam2').val() != $('#halfScoreTeam2').attr('cur')) {
        json['T2half'] = parseInt($('#halfScoreTeam2').val());
    }
    if ($('#endTeam1').val() != $('#endTeam1').attr('cur')) {
        json['T1final'] = parseInt($('#endTeam1').val());
    }
    if ($('#endTeam2').val() != $('#endTeam2').attr('cur')) {
        json['T2final'] = parseInt($('#endTeam2').val());
    }
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
            if (this.readyState == 4 && this.status == 200) {
                $('#save-game-edits').removeClass('btn-primary').addClass('btn-success');
                $('#game-edit-modal').attr('dirty', 1)
            }
        };
        xhttp.open("POST", "edit_game", true)
        // xhttp.setRequestHeader("Content-type", "application/json");
        xhttp.send(getChangedGameJson());
    }
});

$('#game-edit-modal').on('hide.bs.modal', function(e) {
    if ($(this).attr('dirty') == 1) {
        location.reload();
    }
});