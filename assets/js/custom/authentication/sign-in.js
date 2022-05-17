//"use strict";var KTSigninGeneral=function(){var t,e,i;return{init:function(){t=document.querySelector("#kt_sign_in_form"),e=document.querySelector("#kt_sign_in_submit"),i=FormValidation.formValidation(t,{fields:{password:{validators:{notEmpty:{message:"The password is required"}}}},plugins:{trigger:new FormValidation.plugins.Trigger,bootstrap:new FormValidation.plugins.Bootstrap5({rowSelector:".fv-row"})}}),e.addEventListener("click",(function(n){n.preventDefault(),i.validate().then((function(i){"Valid"==i?(e.setAttribute("data-kt-indicator","on"),e.disabled=!0,setTimeout((function(){e.removeAttribute("data-kt-indicator"),e.disabled=!1}),2e3)):Swal.fire({text:"Sorry, looks like there are some errors detected, please try again.",icon:"error",buttonsStyling:!1,confirmButtonText:"Ok, got it!",customClass:{confirmButton:"btn btn-primary"}})}))}))}}}();KTUtil.onDOMContentLoaded((function(){KTSigninGeneral.init()}));

"use strict";
var KTSigninGeneral = (function () {
    var t, e, i;
    return {
        init: function () {
            (t = document.querySelector("#kt_sign_in_form")),
                (e = document.querySelector("#kt_sign_in_submit")),
                (i = FormValidation.formValidation(t, {
                    fields: {
                        password: { validators: { notEmpty: { message: "The password is required" } } },
                    },
                    plugins: { trigger: new FormValidation.plugins.Trigger(), bootstrap: new FormValidation.plugins.Bootstrap5({ rowSelector: ".fv-row" }) },
                })),
                e.addEventListener("click", function (n) {
                    n.preventDefault(),
                        i.validate().then(function (i) {
                            "Valid" == i
                                ? validatePassword(t.querySelector('[name="password"]').value)
                                : function () { };
                        });
                });
        },
    };

    function validatePassword(password) {
        var data = { "password": password };
        $.ajax({
            type: "POST",
            url: "/log-in",
            data: JSON.stringify(data),
            success: function (result) {
                window.location.href = "/";
            },
            error: function (error) {
                Swal.fire({
                    text: "Password not valid",
                    icon: "error",
                    buttonsStyling: !1,
                    confirmButtonText: "Ok, got it!",
                    customClass: { confirmButton: "btn btn-primary" },
                });
            }
            //dataType: ""
        });
    }
})();
KTUtil.onDOMContentLoaded(function () {
    KTSigninGeneral.init();
});
