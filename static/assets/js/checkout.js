$(document).ready(function(){

    
    (function($) {
    // validate contactForm form
    $(function() {
        $('#checkoutForm').validate({
            rules: {
                nameF: {
                    required: true,
                    minlength: 2
                },
                nameL: {
                    required: true,
                    minlength: 1
                },
                company: {
                    required: false,
                },
                email: {
                    required: true,
                    email: true
                },
                number: {
                    required: true,
                },
                add1: {
                    required: true,
                },
                add2: {
                    required: false,
                },
                city: {
                    required: true,
                },
                zip: {
                    required: false,
                },
                country: {
                    required: true,
                }

            },
            messages: {
                nameF: {
                    required: "Please enter a name",
                    minlength: "Your name must consist of at least 2 characters"
                },
                nameL: {
                    required: "Please enter a name",
                },
                email: {
                    required: "Please enter an email"
                },
                number: {
                    required: "Please enter a phone number",
                },
                add1: {
                    required: "Please enter an address"
                },
                city: {
                    required: "Please enter a Town/City"
                },
                country: {
                    required: "Please select a country"
                }
            },
        })
    })
        
 })(jQuery)
})