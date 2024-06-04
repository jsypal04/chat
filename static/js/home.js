let width = window.innerWidth;

if (width <= 800) {
    let conversations = document.getElementById('conversations');
    conversations.style.width = width - 50;
}

async function fetchConvo() {
    const data = await fetch('http://localhost/conversation/')
}