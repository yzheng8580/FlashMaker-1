$(function() {
    $("#debug").html("Text");
});

function enterText() {
    $("#textbox").keypress(function(event) {
        if (event.which == 13) {
            var word = $("#textbox").val();
            $("#textbox").val("");
            var prevState = $("#debug").html() + "<br>";
            $("#debug").html(prevState + word);
        }
    });
}
