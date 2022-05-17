"use strict";
var KTSaveSettings = (function () {
    var e, t, r, o,
        s = function () {
            return 100 === o.getScore();
        };
    return {
        init: function () {
            (e = document.querySelector("#kt_setting_form")),
                (t = document.querySelector("#kt_setting_form_button")),
                //(o = KTPasswordMeter.getInstance(e.querySelector('[data-kt-password-meter="true"]'))),
                (r = FormValidation.formValidation(e, {
                    fields: {
                        "server": {
                            validators: {
                                notEmpty: { message: "The server is required" },
                                callback: {
                                    message: "Please enter valid server",
                                    callback: function (e) {
                                        if (e.value.length > 0) return true;
                                    },
                                },
                            },
                        },
                        "key": {
                            validators: {
                                notEmpty: { message: "The key is required" },
                                callback: {
                                    message: "Please enter valid key",
                                    callback: function (e) {
                                        if (e.value.length > 0) return true;
                                    },
                                },
                            },
                        }
                    },
                    plugins: { trigger: new FormValidation.plugins.Trigger({ event: { newPassword: !1 } }), bootstrap: new FormValidation.plugins.Bootstrap5({ rowSelector: ".fv-row", eleInvalidClass: "", eleValidClass: "" }) },
                })),
                t.addEventListener("click", function (s) {
                    s.preventDefault(),
                        r.validate().then(function (r) {
                            "Valid" == r
                                ? (t.setAttribute("data-kt-indicator", "on"),
                                    (t.disabled = !0),
                                    saveSettings()
                                )
                                : function () { }
                        });
                }),
                $("#validate").on('change', function() {
                    console.log($(this).is(':checked'));
                })
                /*e.querySelector('input[name="validate"]').addEventListener("input", function () {
                    console.log();
                })*/;
        },
    };

    function saveSettings() {
        let validate = $("#validate").is(':checked') == false ? 0 : 1;
        var data = { "id": Number(e.querySelector('[name="IdSet"]').value), "server": e.querySelector('[name="server"]').value, "key": e.querySelector('[name="key"]').value, "validateCertificate": validate };
        console.log(data);
        $.ajax({
            type: "POST",
            url: "/save-settings",
            data: JSON.stringify(data),
            success: function (result) {
                t.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: "You have successfully save settings!",
                    icon: "success",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" }
                }).then(function(){
                    //window.location.href = "/"
                });
                t.disabled = !1
            },
            error: function (error) {
                t.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: "Sorry, there was an error. Try again!",
                    icon: "error",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" },
                });
                t.disabled = !1;
            }
        });
    }
})();
KTUtil.onDOMContentLoaded(function () {
    KTSaveSettings.init();
});