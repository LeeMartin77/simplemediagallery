// const button = document.querySelector("#dir-toggle-button");

// button.addEventListener("click", (event) => {
//   button.textContent = `Click count: ${event.detail}`;
//   console.log(event)
// });

function toggleDirectories() {
    const button = document.getElementById("dir-toggle-button");
    const directoryContainer = document.getElementById('directories');

    isHidden = directoryContainer.style.visibility == 'hidden';
    directoryContainer.style.visibility = isHidden ? 'visible' : 'hidden';
    directoryContainer.style.maxHeight = isHidden ? '12em' : '0';
    directoryContainer.style.padding = isHidden ? '1em' : '0';
    button.textContent = isHidden ? 'Directories' : 'Collapse';
}