$(document).ready(function () {
    loadPage(1, null);
    var typedMain = new Typed('#typed', {
        stringsElement: "#typedElements",
        typeSpeed: 40,
        backSpeed: 10,
        backDelay: 700,
        loop: false,
        onComplete: function (self) {
            $("#typed").append("<strong><i> ⇇</i></strong>")
            $("#typed").prepend("<strong><i>⇉ </i></strong>")
        }
    });

    $("#typed").click(function () {
        typedMain.reset();
    });

    if (typedMain !== undefined) {
        $("#typedElements").hide();
    }

    $("#OlderPosts").click(function () {
        loadPage($("#OlderPosts").attr("data-next-page"), $("#OlderPosts").parent());
    });

});

var loadingPages = false;
function loadPage(pageNumber, removeElement) {
    if (loadingPages) {
        console.log("ALREADY LOADING PAGES!");
        return;
    }

    if (!Number.isInteger(pageNumber)) {
        console.log("NOT AN INTEGER!");
        return;
    }

    $("#posts").append("<div class='loader' id='pageLoading'></div>");
    $.get("/post_preview/" + pageNumber, function () {})
        .done(function (data) {
            $("#posts").append(data);
            if (removeElement !== undefined && removeElement !== null) {
                removeElement.remove();
            }
        })
        .fail(function (jqXHR, textStatus, error) {
            console.log(textStatus, error);
        })
        .always(function () {
            loadingPages = false;
            $("#pageLoading").remove();
        });
}