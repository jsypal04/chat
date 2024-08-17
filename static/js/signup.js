function emailIsValid(email) {
    let atIndex = email.indexOf("@");
    if (atIndex == -1) { 
        return false; 
    }

    let domain = email.slice(atIndex + 1);
    if (domain != "cua.edu") {
        return false;
    }

    return true;
}

function validateEmail() {
    let email = document.getElementsByName("email");
    if (!emailIsValid(email[0].value) && email[0].value != "") {
        document.getElementById("invalid-email").innerText = "You need a cua.edu email to sign up";
    }
    else {
        document.getElementById("invalid-email").innerText = "";
    }
}

function comparePasswords() {
    let password = document.getElementsByName("password")[0].value;
    let confirmation = document.getElementsByName("confirmPassword")[0].value;
    if (password != confirmation && confirmation != "") {
        document.getElementById("pw-discrepency").innerText = "Passwords must match";
    }
    else {
        document.getElementById("pw-discrepency").innerText = "";
    }
}

function validateForm() {
    let data = document.getElementsByClassName("text-input");
    let password;
    for (let i = 0; i < data.length; i++) {
        if (data[i].name == "email" && !emailIsValid(data[i].value)) {
            event.preventDefault();
        }
        else if (data[i].name == "password") {
            password = data[i].value;
        }
        else if (data[i].name == "confirm-password" && password != data[i].value) {
            event.preventDefault();
        }
    }
}