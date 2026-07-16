/* ============================================================
   AMJ HUB — Portfolio Page JavaScript
   Handles: category filter for portfolio case study cards
   ============================================================ */

(function () {
  'use strict';

  const filterBtns    = document.querySelectorAll('.portfolio-filter-bar .gallery-filter');
  const portfolioCards = document.querySelectorAll('#portfolio-cards .portfolio-card');

  if (!filterBtns.length) return;

  filterBtns.forEach(function (btn) {
    btn.addEventListener('click', function () {
      const filter = this.dataset.filter;

      // Update active state on buttons
      filterBtns.forEach(function (b) {
        b.classList.remove('active');
        b.setAttribute('aria-selected', 'false');
      });
      this.classList.add('active');
      this.setAttribute('aria-selected', 'true');

      // Show / hide portfolio cards
      portfolioCards.forEach(function (card) {
        // data-category can contain multiple space-separated values e.g. "video brand"
        const cats = (card.dataset.category || '').split(' ');
        if (filter === 'all' || cats.includes(filter)) {
          card.classList.remove('hidden');
        } else {
          card.classList.add('hidden');
        }
      });
    });
  });

})();
