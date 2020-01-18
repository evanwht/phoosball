$("#alert-box").hide();

(".goal-rating").on("change", function(event) {
    if (event.val() > 5 || event.val() < 0) {
        // hightlight red and show pop up/ tool tip saying value needs to be in range
    }
});
