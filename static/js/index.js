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
        document.getElementById('chat').style.display = "flex";
        document.getElementById('conversations').style.width = "350px";
        document.getElementById('chat').style.width = (width - 350) + "px";
    }

    // set the chat-form to the proper position and size
    let chatForm = document.getElementById('chat-form');
    chatForm.style.left = Math.floor(0.05 * (width - 400) + 400) + "px";
    chatForm.style.width = Math.floor(0.9 * (width - 400)) + "px";
}

function showNewChatFlag() {
    setTimeout(() => {
        document.getElementById("new-chat-flag").style.display = "flex";
    }, 500);
}

function hideNewChatFlag() {
    document.getElementById("new-chat-flag").style.display = "none";
}

/**
 * A function to show the user select dropdown
 */
function showUserDropdown() {
    let dropdown = document.getElementById("select-user-dropdown");
    dropdown.style.display = "block";
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
 * A function to show the modal to create a new conversation.
 * Also gets a list of users for the user to select from
 */
async function showNewConvoModal() {
    let screen = document.getElementById("new-convo-screen");
    screen.style.display = "flex";

    // get a list of users
    let emails = await fetch("http://localhost/get-users")
        .then((httpRes) => {
            if (httpRes.ok) {
                return httpRes.json();
            }
            console.log("Request failed with status code " + httpRes.status);
        }).catch((err) => { console.log(err); });

    // render a dropdown selection
    let keys = Object.keys(emails);
    let dropdown = document.getElementById("select-user-dropdown");
    for (let i = 0; i < keys.length; i++) {
        let option = document.createElement("p");
        option.innerText = emails[keys[i]] + " (" + keys[i] + ")";
        option.classList = "user-option";
        option.addEventListener("click", () => {
            document.getElementById("select-user-dropdown").style.display = "none";
            document.getElementById("new-recipient").value = keys[i];
        });
        dropdown.appendChild(option);
    }
    
}

/**
 * A function to close the modal to create a new conversation
 */
function closeNewConvoModal() {
    let screen = document.getElementById("new-convo-screen");
    screen.style.display = "none";
    document.getElementById("select-user-dropdown").style.display = "none";
    document.getElementById("new-recipient").value = "";
    let users = document.getElementsByClassName("user-option");
    while (users.length > 0) {
        users[0].remove();
    }
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
    while (myMessages.length > 0) {
        myMessages.item(0).remove();
    }    
    while (otherMessages.length > 0) {
        otherMessages.item(0).remove();
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
    
    let parent = document.getElementById("chat-messages");
    for (let i = 0; i < res.length; i++) {
        let message = document.createElement('div');
        if (res[i].sender == 'Me') {
            message.className = "my-bubble";
        } else {
            message.className = "other-bubble";
        }
        message.innerText = res[i].content;
        parent.appendChild(message);
    }

    currentOpenConvo = ID;
    document.getElementById(ID).style.backgroundColor = "#E6E6E6";
}

/**
 * A function to to send a message to the server
 */
async function sendMessage() {
    if (currentOpenConvo == null) {
        console.log("You cannot send a message without a conversation open.");
        return;
    }
    let date = new Date();
    let time = date.getTime();
    let content = document.getElementById('entry');
    let res = await fetch("http://localhost/send", {
        method: 'POST',
        body: JSON.stringify({
            id: time,
            convoID: parseInt(currentOpenConvo),
            content: content.value
        })
    }).then((httpRes) => {
        if (httpRes.ok) {
            return httpRes.json();
        }
        console.log("Request failed with status code: " + httpRes.status);
    }).catch((err) => { console.log(err); });

    let newMessage = document.createElement("div");
    newMessage.className = "my-bubble";
    newMessage.innerText = res.content;
    document.getElementById('chat-messages').appendChild(newMessage);

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