//DEMONSTRATING HOW JAVASCRIPT MODULES WORK
function prompt() {
    //CREATING TOAST ALERT
    let toast = function (c) {
        //parameter "c" will be overwritten by the values in the const assigned to "c"
        const {
            msg = "Signed in successfully", //default values if not specified
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            position: position,
            title: msg,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })
        Toast.fire({}) //firing the sweetalert
    }
    //*****************************************************

    //CREATING SUCCESS ALERT
    let success = function (c) {
        const {
            icon = "success",
            msg = "",
            title = "",
            footer = "",

        } = c;
        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
        })
    }

    //CREATING THE ERROR ALERT
    let error = function (c) {
        const {
            icon = "error",
            msg = "",
            title = "",
            footer = "",

        } = c;
        Swal.fire({
            icon: icon,
            title: title,
            text: msg,
            footer: footer,
        })
    }

    //BELOW CODE IS USED ON MAJORS AND GENERALS PAGE AS A MODAL
    //CREATING MULTIPLE INPUTS: FOR ARRIVAL AND DEPARTURE DATES FOR USERS, IT USES SWEETALERT IN MAJORS AND GENERAL PAGES
    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c; //code explain: when the function is call we pass it a message and a title which then set the values we have passed to the var "c"

        /*whatever is put into the form and do something with it*/
        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            //before it opens it is going to show the datepicker
            willOpen: () => {
                //if willOpen is specified when I make the call to this custom function then execute willOpen
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },

            preConfirm: () => {
                //shows what was entered by the user
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
            //after it is opens it is going to remove the disabled attribute on the input element
            didOpen: () => {
                if (c.didOpen !== undefined){
                    c.didOpen();
                }
            }
        })
        //just show what is entered but not what I want
        /*if (formValues) {
          Swal.fire(JSON.stringify(formValues))
        }*/

        //WORKING ON THE AJAX REQUEST
        //checking to see if we actually have a result
        if (result){
            //if they actually do not click on the cancel button then want to do something
            if (result.dismiss !== Swal.DismissReason.cancel){
                //(result.value !== ""): if there is result value
                if (result.value !== ""){
                    //I want to call a callback, want to execute some js back on the clicked page thus what happen when someone fill a form
                    //if there is a callback then want to do something
                    if (c.callback !== undefined){
                        //calling the callback and pass it the result(what they entered or what I got back)
                        c.callback(result);
                    }
                } else {
                    //if the result value is exactly equal to nothing, I want to do nothing
                    c.callback(false);
                }
            } else {
                //(result.dismiss !== Swal.DismissReason.cancel): if the result is dismissed, if they hit the cancel button then do not do anything
                c.callback(false);
            }
        } //CODE EXPLANATION: it allows me to process code after the Swal dialog box is closed thus after they hit the submit button

    }

    return {
        toast: toast, //if there is a request for toast, return the variable holding the function toast
        success: success,
        error: error,
        custom: custom,
    }

}