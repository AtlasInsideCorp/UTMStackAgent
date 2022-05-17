"use strict";
var KTUpdateModule = (function () {
    var e, t, r, o,
        s = function () {
            return 100 === o.getScore();
        };
    return {
        init: function () {
            (e = document.querySelectorAll(".kt_module_form")),
                (t = document.querySelectorAll(".kt_module_form_button")),
                t.forEach(x => {
                    x.addEventListener("click", function (s) {
                        s.preventDefault(),
                            (x.setAttribute("data-kt-indicator", "on"),
                                (x.disabled = !0),
                                updateModule(x)
                            )
                    });
                });

        },
    };

    function updateModule(button) {
        let enable = Number(button.getAttribute("kt-enable")) == 0 ? 1 : 0; 
        var data = { "id": Number(button.getAttribute("id")), "name": "", "image": "", "enable": enable };
        console.log(data);
        $.ajax({
            type: "POST",
            url: "/edit-module",
            data: JSON.stringify(data),
            success: function (result) {
                button.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: "You have successfully edited module!",
                    icon: "success",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" }
                }).then(function () {
                    window.location.href = "/"
                });
                button.disabled = !1
            },
            error: function (error) {
                button.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: "Sorry, there was an error. Try again!",
                    icon: "error",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" },
                });
                button.disabled = !1;
            }
        });
    }
})();
KTUtil.onDOMContentLoaded(function () {
    KTUpdateModule.init();
});