function returnToDefaultPosition(square) {
    let defaultX = square.getAttribute('data-default-x-px');
    let defaultY = square.getAttribute('data-default-y-px');
    square.style.left = `${defaultX}px`;
    square.style.top = `${defaultY}px`;
    square.removeAttribute('data-X-pox');
    square.removeAttribute('data-Y-pos');
    document.getElementById("availble-characters").appendChild(square);
}

function returnAllSquaresToDefaultPosition() {
    let draggableSquares = document.querySelectorAll('.draggable-square')
    draggableSquares.forEach(square => {
        returnToDefaultPosition(square);
    });
}

document.getElementById('rotate-left').addEventListener('click', () => {
    returnAllSquaresToDefaultPosition();
})

document.getElementById('rotate-right').addEventListener('click', () => {
    returnAllSquaresToDefaultPosition();
})

document.getElementById('test-functions').addEventListener('click', () => {
    returnAllSquaresToDefaultPosition();
})

document.getElementById('make-play').addEventListener('click', () => {
    returnAllSquaresToDefaultPosition();
})