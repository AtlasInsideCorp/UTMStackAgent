"use strict";
var KTSignoutGeneral = (function () {
    var t, e;
    return {
        init: function () {
            (t = document.querySelector("#kt_header_user_menu_toggle")),
                (e = document.querySelector("#kt_header_user_menu_toggle_logout")),
                e.addEventListener("click", function (n) {
                    n.preventDefault(),
                    signOut()
                });
        },
    };

    function signOut() {
        $.ajax({
            type: "POST",
            url: "/log-out",
            data: {},
            success: function (result) {
                window.location.href = "/sign-in";
            },
            error: function (error) {
            }
        });
    }
})();
KTUtil.onDOMContentLoaded(function () {
    KTSignoutGeneral.init();
});
