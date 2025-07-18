document.addEventListener('DOMContentLoaded', () => {
    let isDragging = false;
    let startX, startY;
    let activeSquare = null;

    function initializeDraggableSquares() {
        let draggableSquares = document.querySelectorAll('.draggable-square:not([data-initialized])');
        
        draggableSquares.forEach((square, index) => {
            if (!square.hasAttribute('data-default-x-px')) {
                square.setAttribute('data-default-x-px', 20 + index * 60);
                square.setAttribute('data-default-y-px', 20);
                square.style.position = 'absolute';
                square.style.left = `${square.getAttribute('data-default-x-px')}px`;
                square.style.top = `${square.getAttribute('data-default-y-px')}px`;
            }
            
            square.addEventListener('mousedown', dragStart);
            square.setAttribute('data-initialized', 'true');
        });
    }

    initializeDraggableSquares();

    document.body.addEventListener('htmx:wsAfterMessage', function(event) {
        initializeDraggableSquares();
    });

    document.addEventListener('mouseup', dragEnd);
    document.addEventListener('mousemove', drag);

    function dragStart(e) {
        activeSquare = e.target;
        startX = e.clientX;
        startY = e.clientY;
        isDragging = true;
    }

    function dragEnd(e) {
        if (!isDragging) return;

        isDragging = false;
        let targetFace = getActiveFace(gridContainers);
        if (isWithinOuterCube(e.clientX, e.clientY, targetFace) && targetFace) {
            let faceRect = targetFace.getBoundingClientRect();
            let gridSize = Math.round(faceRect.width / 15);

            let relativeX = Math.round(e.clientX - faceRect.left);
            let relativeY = Math.round(e.clientY - faceRect.top);

            let squareCenterOffsetX = Math.round(activeSquare.offsetWidth / 2);
            let squareCenterOffsetY = Math.round(activeSquare.offsetHeight / 2);

            relativeX -= squareCenterOffsetX;
            relativeY -= squareCenterOffsetY;

            relativeX = Math.round(relativeX / gridSize) * gridSize;
            relativeY = Math.round(relativeY / gridSize) * gridSize;

            relativeX = Math.max(0, Math.min(relativeX, faceRect.width - activeSquare.offsetWidth));
            relativeY = Math.max(0, Math.min(relativeY, faceRect.height - activeSquare.offsetHeight));

            // Obliczanie pozycji w siatce (0-14)
            let gridX = Math.floor(relativeX / gridSize);
            let gridY = Math.floor(relativeY / gridSize);
            
            // Ustawianie atrybutów data-X i data-Y
            activeSquare.setAttribute('data-X-pos', gridX);
            activeSquare.setAttribute('data-Y-pos', gridY);

            activeSquare.style.left = `${relativeX}px`;
            activeSquare.style.top = `${relativeY}px`;

            targetFace.appendChild(activeSquare);
        } else {
            returnToDefaultPosition(activeSquare);
        }
        activeSquare = null;
    }

    function drag(e) {
        if (isDragging) {
            e.preventDefault();
            let dx = e.clientX - startX;
            let dy = e.clientY - startY;

            let newLeft = activeSquare.offsetLeft + dx;
            let newTop = activeSquare.offsetTop + dy;

            activeSquare.style.left = `${newLeft}px`;
            activeSquare.style.top = `${newTop}px`;

            startX = e.clientX;
            startY = e.clientY;
        }
    }

    function isWithinOuterCube(x, y, targetFace) {
        if (!targetFace) return false;
        let rect = targetFace.getBoundingClientRect();
        return x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom;
    }
});