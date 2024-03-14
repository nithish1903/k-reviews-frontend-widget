document.addEventListener("DOMContentLoaded", function () {
    var allPTags = document.querySelectorAll('.reviewText');

    allPTags.forEach(function (pTag) {
        var words = pTag.textContent.split(' ');

        if (words.length > 80) {
            var truncatedText = words.slice(0, 80).join(' ');
            pTag.textContent = truncatedText + '...';
        }
    });
});