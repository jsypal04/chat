resize();

function resize() {
    let width = window.innerWidth;

    if (width <= 800) {
        document.getElementById('conversations').style.width = (width - 50) + "px";
        document.getElementById('chat').style.display = "none";
    } else {
        document.getElementById('conversations').style.width = "350px";
        document.getElementById('chat').style.width = (width - 527) + "px";
    }
}

function openConversation(ID) {
    let convo = document.getElementById(ID);
    fetch("/id/" + ID, {
        method: 'POST',
        body: JSON.stringify({
            Id: ID,
            Sender: "Me",
            Receiver: convo.innerText,
        })
    });
}

function sendMessage() {
    let content = document.getElementById('entry');
    fetch("/", {
        method: 'POST',
        body: content.value
    });
    content.value = "";
}