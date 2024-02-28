function showTabKP(tabId, buttonId,tab,activeTab,btn,activeBtn) {
    // Hide all tabs
    var tabs = document.getElementsByClassName(tab);
    for (var i = 0; i < tabs.length; i++) {
    tabs[i].classList.remove(activeTab);
    }

    // Show the selected tab
    document.getElementById(tabId).classList.add(activeTab);

    // Remove "active-button" class from all buttons
    var buttons = document.getElementsByClassName(btn);
    for (var j = 0; j < buttons.length; j++) {
        buttons[j].classList.remove(activeBtn);
    }

    // Add "active-button" class to the clicked button
    document.getElementById(buttonId).classList.add(activeBtn);
}