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
    chatForm.style.left = Math.floor(0.05 * (width - 400) + 400) + "px";;
    chatForm.style.width = Math.floor(0.9 * (width - 400)) + "px";
}

/**
 * A function to fetch the the messages associated with a conversation id and render them in the chat area
 * @param {number} ID - the id number of a conversation 
 */
async function openConversation(ID) {
    let myMessages = document.getElementsByClassName('my-bubble');
    let otherMessages = document.getElementsByClassName('other-bubble');
    for (let i = 0; i < myMessages.length; i++) {
        myMessages[i].remove();
    }
    console.log(myMessages);
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
    
    let parent = document.getElementById('chat-messages');
    for (let i = 0; i < res.content.length; i++) {
        let message = document.createElement('div');
        if (res.content[i].sender == 'Me') {
            message.className = "my-bubble";
        } else {
            message.className = "other-bubble";
        }
        message.style.top = (75 * i) + "px";
        message.innerText = res.content[i].content;
        parent.appendChild(message);
    }
}

/**
 * A function to to send a message to the server
 */
function sendMessage() {
    let date = new Date();
    let time = date.getTime();
    let content = document.getElementById('entry');
    fetch("http://localhost/send", {
        method: 'POST',
        body: JSON.stringify({
            id: time,
            sender: "Me",
            receiver: "Bob",
            content: content.value
        })
    });
    content.value = "";
}