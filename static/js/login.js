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
        document.getElementById("invalid-email").innerText = "You need a cua.edu email to sign in";
    }
    else {
        document.getElementById("invalid-email").innerText = "";
    }
}

function validateForm() {
    let email = document.getElementsByName("email");
    if (!emailIsValid(email[0].value)) {
        event.preventDefault();
    }
}