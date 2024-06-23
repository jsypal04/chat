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