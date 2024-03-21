document.addEventListener('DOMContentLoaded', function () {
    // Get the buttons and tabs
    const siteReviewsBtn = document.getElementById('kpr-SiteReviews-btn');
    const productReviewsBtn = document.getElementById('kpr-ProductReviews-btn');
    
    const siteReviewsTab = document.getElementById('kpr-SiteReviews-tab');
    const productReviewsTab = document.getElementById('kpr-ProductReviews-tab');

    // Function to toggle between tabs
    function toggleTabs(activeBtn, activeTab, inactiveBtn, inactiveTab) {
        activeBtn.classList.add('kpSP-active-button');
        inactiveBtn.classList.remove('kpSP-active-button');

        activeTab.classList.remove('hidden');
        inactiveTab.classList.add('hidden');
    }

    // Add click event listeners to buttons
    siteReviewsBtn.addEventListener('click', function () {
        toggleTabs(siteReviewsBtn, siteReviewsTab, productReviewsBtn, productReviewsTab);
    });

    productReviewsBtn.addEventListener('click', function () {
        toggleTabs(productReviewsBtn, productReviewsTab, siteReviewsBtn, siteReviewsTab);
    });
});
