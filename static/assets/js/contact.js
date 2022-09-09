$(document).ready(function(){

    
    (function($) {
    // validate contactForm form

    $(function() {
        $('#contactForm').validate({
            rules: {
                
                name: {
                    required: true,
                    minlength: 2
                },
                subject: {
                    required: false,
                    minlength: 0
                },
                number: {
                    required: false,
                    minlength: 5
                },
                email: {
                    required: true,
                    email: true
                },
                message: {
                    required: true,
                    minlength: 10
                }
            },
            messages: {
                name: {
                    required: "Please enter a name",
                    minlength: "Your name must consist of at least 2 characters"
                },
                email: {
                    required: "Please enter an email"
                },
                message: {
                    required: "Please enter a message",
                    minlength: "The message needs to be at least 10 characters long."
                }
            },
        })
    })
        
 })(jQuery)
})