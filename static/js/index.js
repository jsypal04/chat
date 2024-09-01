let currentOpenConvo = null;

resize();

/**
 * Function to resize the window based on the window width
 */
function resize() {
    let width = window.innerWidth;

    if (width <= 800) {
        document.getElementById('conversations').style.width = (width - 50) + "px";
        document.getElementById('chat').style.display = "none";
    } else {
        document.getElementById('chat').style.display = "block";
        document.getElementById('conversations').style.width = "350px";
        document.getElementById('chat').style.width = (width - 400) + "px";
    }

    // set the chat-form to the proper position and size
    let chatForm = document.getElementById('chat-form');
    chatForm.style.left = Math.floor(0.05 * (width - 400) + 400) + "px";
    chatForm.style.width = Math.floor(0.9 * (width - 400)) + "px";
}

/**
 * A function to show and close the profile options for a user
 */
function showProfileOptions() {
    let profileOptions = document.getElementById("profile-options");
    if (profileOptions.style.display == "none") {
        profileOptions.style.display = "block";
    }
    else if (profileOptions.style.display == "block") {
        profileOptions.style.display = "none";
    }
}

/**
 * A function to show the modal to create a new conversation
 */
function showNewConvoModal() {
    let screen = document.getElementById("new-convo-screen");
    screen.style.display = "flex";
}

/**
 * A function to close the modal to create a new conversation
 */
function closeNewConvoModal() {
    let screen = document.getElementById("new-convo-screen");
    screen.style.display = "none";
}

/**
 * A function to fetch the the messages associated with a conversation id and render them in the chat area
 * @param {number} ID - the id number of a conversation 
 */
async function openConversation(ID) {
    if (currentOpenConvo != null) {
        document.getElementById(currentOpenConvo).style.backgroundColor = "#F0F0F0";
    }
    document.getElementById("chat-form").style.display = "block";
    let myMessages = document.getElementsByClassName('my-bubble');
    let otherMessages = document.getElementsByClassName('other-bubble');
    for (let i = 0; i < myMessages.length; i++) {
        myMessages[i].remove();
    }
    for (let i = 0; i < otherMessages.length; i++) {
        otherMessages[i].remove();
    }

    let res = await fetch("http://localhost/id/" + ID, {method: 'GET'})
        .then((httpRes) => {
            if (httpRes.ok) {
                return httpRes.json();
            }
            console.log("Http request failed with status code " + httpRes.status);
        })
        .catch((err) => { console.log(err); });
    
    // if the response is null (ie, empty) simply set the current open convo and return
    if (res == null) {
        currentOpenConvo = ID;
        document.getElementById(ID).style.backgroundColor = "#E6E6E6";
        return;
    }

    let parent = document.getElementById('chat-messages');
    for (let i = 0; i < res.length; i++) {
        let message = document.createElement('div');
        if (res[i].sender == 'Me') {
            message.className = "my-bubble";
        } else {
            message.className = "other-bubble";
        }
        message.style.top = (75 * i) + "px";
        message.innerText = res[i].content;
        parent.appendChild(message);
    }

    currentOpenConvo = ID;
    document.getElementById(ID).style.backgroundColor = "#E6E6E6";
}

/**
 * A function to to send a message to the server
 */
function sendMessage() {
    if (currentOpenConvo == null) {
        console.log("You cannot send a message without a conversation open.");
        return;
    }
    console.log(currentOpenConvo);
    let date = new Date();
    let time = date.getTime();
    let content = document.getElementById('entry');
    fetch("http://localhost/send", {
        method: 'POST',
        body: JSON.stringify({
            id: time,
            convoID: parseInt(currentOpenConvo),
            content: content.value
        })
    });
    content.value = "";
}

/**
 * A function to create a new conversation
 */
function createConvo() {
    // get the email of the receiver
    let recipientAddr = document.getElementById("new-recipient");
    // close the modal
    document.getElementById("new-convo-screen").style.display = "none";
    // get the rest of the data
    let date = new Date();
    let time = date.getTime();
    fetch("http://localhost/new-convo", {
        method: 'POST',
        body: JSON.stringify({
            id: time,
            receiver: recipientAddr.value
        })
    });
    recipientAddr.value = "";
}