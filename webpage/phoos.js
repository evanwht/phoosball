var players =  JSON.parse($("#choosePlayers").attr("players"));
var goal_types =  JSON.parse($("#goal-sect").attr("goal-types"));

var option = function(text) {
  return '<option value="' + text + '">' + text + '</option>';
};

var updatePlayers = function() {
  $(".goal-scorer").empty().append(option(""));
  if ($("#player1").val() != "") {
    $(".goal-scorer").append(option($("#player1").val()));
  }
  if ($("#player2").val() != "") {
    $(".goal-scorer").append(option($("#player2").val()));
  }
  if ($("#player3").val() != "") {
    $(".goal-scorer").append(option($("#player3").val()));
  }
  if ($("#player4").val() != "") {
    $(".goal-scorer").append(option($("#player4").val()));
  }
}

goal_types.forEach(function(type, index) {
  $(".goal-type").append(option(type));
});

players.forEach(function(player, index) {
  $("#player1").append(option(player));
  $("#player2").append(option(player));
  $("#player3").append(option(player));
  $("#player4").append(option(player));
});

$("#player1").on("change", function(event) {
  updatePlayers();
});
$("#player2").on("change", function(event) {
  updatePlayers();
});
$("#player3").on("change", function(event) {
  updatePlayers();
});
$("#player4").on("change", function(event) {
  updatePlayers();
});
