let currentRotation = 0;

document.getElementById('rotate-left').addEventListener('click', () => {
    currentRotation -= 90;
    document.querySelector('#outer-cube').style.transform = `rotateY(${currentRotation}deg)`;
});

document.getElementById('rotate-right').addEventListener('click', () => {
    currentRotation += 90;
    document.querySelector('#outer-cube').style.transform = `rotateY(${currentRotation}deg)`;
});

function getNormalizedRotation() {
    return ((currentRotation % 360) + 360) % 360;
}

function getActiveFace(gridContainers) {
    let normalizedRotation = getNormalizedRotation();

    switch(normalizedRotation) {
        case 0:   return gridContainers[0]; // north 
        case 90:  return gridContainers[1]; // east
        case 180: return gridContainers[2]; // south
        case 270: return gridContainers[3]; // west
        default:  throw `wrong normalizedRotation ${normalizedRotation};`;
    }
}
