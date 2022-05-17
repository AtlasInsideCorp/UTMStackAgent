"use strict";
var KTChangePassword = (function () {
    var e, t, r, o,
        s = function () {
            return 100 === o.getScore();
        };
    return {
        init: function () {
            (e = document.querySelector("#kt_change_password_form")),
                (t = document.querySelector("#kt_change_password_submit")),
                (o = KTPasswordMeter.getInstance(e.querySelector('[data-kt-password-meter="true"]'))),
                (r = FormValidation.formValidation(e, {
                    fields: {
                        "current-password": {
                            validators: {
                                notEmpty: { message: "The password is required" },
                                callback: {
                                    message: "Please enter valid password",
                                    callback: function (e) {
                                        if (e.value.length > 0) return true;
                                    },
                                },
                            },
                        },
                        "new-password": {
                            validators: {
                                notEmpty: { message: "The new password is required" },
                                callback: {
                                    message: "Please enter valid password",
                                    callback: function (e) {
                                        if (e.value.length > 0) return s();
                                    },
                                },
                            },
                        },
                        "confirm-password": {
                            validators: {
                                notEmpty: { message: "The password confirmation is required" },
                                identical: {
                                    compare: function () {
                                        return e.querySelector('[name="new-password"]').value;-p
                                    },
                                    message: "The new password and its confirm are not the same",
                                },
                            },
                        },
                        toc: { validators: { notEmpty: { message: "You must accept the terms and conditions" } } },
                    },
                    plugins: { trigger: new FormValidation.plugins.Trigger({ event: { newPassword: !1 } }), bootstrap: new FormValidation.plugins.Bootstrap5({ rowSelector: ".fv-row", eleInvalidClass: "", eleValidClass: "" }) },
                })),
                t.addEventListener("click", function (s) {
                    s.preventDefault(),
                        r.revalidateField("new-password"),
                        r.validate().then(function (r) {
                            "Valid" == r
                                ? (t.setAttribute("data-kt-indicator", "on"),
                                    (t.disabled = !0),
                                    changePassword(e.querySelector('[name="current-password"]').value, e.querySelector('[name="new-password"]').value)
                                )
                                : function () { }
                        });
                }),
                e.querySelector('input[name="new-password"]').addEventListener("input", function () {
                    this.value.length > 0 && r.updateFieldStatus("new-password", "NotValidated");
                });
        },
    };

    function changePassword(password, newPassword) {
        var data = { "password": password, "newPassword": newPassword };
        console.log(data);
        $.ajax({
            type: "POST",
            url: "/change-password",
            data: JSON.stringify(data),
            success: function (result) {
                t.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: "You have successfully change your password!",
                    icon: "success",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" }
                }).then(function(){
                    window.location.href = "/"
                });
                t.disabled = !1
            },
            error: function (error) {
                t.removeAttribute("data-kt-indicator");
                Swal.fire({
                    text: error.responseJSON.error,
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
    KTChangePassword.init();
});