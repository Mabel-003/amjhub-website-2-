/* ============================================================
   AMJ HUB — Main JavaScript
   Handles: sticky nav, mobile menu, scroll reveals, form AJAX
   ============================================================ */

(function () {
  'use strict';

  /* ── Sticky Navigation ─────────────────────────────────── */
  const navbar = document.getElementById('navbar');

  function updateNav() {
    if (window.scrollY > 40) {
      navbar.classList.add('scrolled');
    } else {
      navbar.classList.remove('scrolled');
    }

    // Active section highlighting
    const sections = document.querySelectorAll('section[id]');
    let current = '';
    sections.forEach(section => {
      const sectionTop = section.offsetTop - 120;
      if (window.scrollY >= sectionTop) {
        current = section.getAttribute('id');
      }
    });

    document.querySelectorAll('.nav-links a').forEach(link => {
      link.classList.remove('active');
      const href = link.getAttribute('href');
      const hash = href.includes('#') ? '#' + href.split('#')[1] : '';
      if (hash === '#' + current) {
        link.classList.add('active');
      }
    });
  }

  window.addEventListener('scroll', updateNav, { passive: true });
  updateNav();

  /* ── Mobile Menu Toggle ────────────────────────────────── */
  const hamburger = document.getElementById('hamburger');
  const mobileMenu = document.getElementById('mobile-menu');

  if (hamburger && mobileMenu) {
    hamburger.addEventListener('click', function () {
      const isOpen = mobileMenu.classList.toggle('open');
      hamburger.classList.toggle('open', isOpen);
      document.body.style.overflow = isOpen ? 'hidden' : '';
    });

    // Close on link click
    mobileMenu.querySelectorAll('a').forEach(link => {
      link.addEventListener('click', function () {
        mobileMenu.classList.remove('open');
        hamburger.classList.remove('open');
        document.body.style.overflow = '';
      });
    });
  }

  /* ── Smooth Scroll for Anchor Links ────────────────────── */
  document.querySelectorAll('a[href^="#"], a[href^="/#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      const href = this.getAttribute('href');
      const targetId = href.startsWith('/#') ? href.substring(1) : href;
      if (targetId === '#') return;
      const target = document.querySelector(targetId);
      if (target) {
        e.preventDefault();
        const navHeight = 80;
        const targetPos = target.getBoundingClientRect().top + window.scrollY - navHeight;
        window.scrollTo({ top: targetPos, behavior: 'smooth' });
      }
    });
  });

  /* ── Scroll Reveal Animation ───────────────────────────── */
  const revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          revealObserver.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1, rootMargin: '0px 0px -40px 0px' }
  );

  document.querySelectorAll('.reveal').forEach(el => {
    revealObserver.observe(el);
  });

  /* ── Contact Form AJAX Submission ──────────────────────── */
  const contactForm = document.getElementById('contact-form');
  const formMessage = document.getElementById('form-message');
  const submitBtn = document.getElementById('form-submit-btn');

  if (contactForm) {
    contactForm.addEventListener('submit', async function (e) {
      e.preventDefault();

      // UI: loading state
      const originalText = submitBtn.innerHTML;
      submitBtn.disabled = true;
      submitBtn.innerHTML = '<span class="btn-spinner"></span> Sending…';
      formMessage.className = '';
      formMessage.style.display = 'none';
      formMessage.textContent = '';

      try {
        const formData = new FormData(contactForm);

        const response = await fetch('/contact', {
          method: 'POST',
          body: new URLSearchParams(formData),
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
        });

        const data = await response.json();

        formMessage.style.display = 'block';

        if (data.success) {
          formMessage.className = 'success';
          formMessage.textContent = '✓ ' + data.message;
          contactForm.reset();
        } else {
          formMessage.className = 'error';
          formMessage.textContent = '✕ ' + (data.error || 'Something went wrong. Please try again.');
        }
      } catch (err) {
        formMessage.style.display = 'block';
        formMessage.className = 'error';
        formMessage.textContent = '✕ Network error. Please check your connection and try again.';
        console.error('Form submission error:', err);
      } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;

        // Scroll form message into view
        formMessage.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
      }
    });
  }

  /* ── Gallery Filter Tabs (homepage) ────────────────────── */
  const galleryFilters = document.querySelectorAll('#gallery .gallery-filter');
  const galleryItems   = document.querySelectorAll('#gallery-grid .gallery-item');

  if (galleryFilters.length) {
    galleryFilters.forEach(btn => {
      btn.addEventListener('click', function () {
        const filter = this.dataset.filter;

        // Update active tab
        galleryFilters.forEach(b => {
          b.classList.remove('active');
          b.setAttribute('aria-selected', 'false');
        });
        this.classList.add('active');
        this.setAttribute('aria-selected', 'true');

        // Filter items
        galleryItems.forEach(item => {
          const category = item.dataset.category || '';
          if (filter === 'all' || category === filter) {
            item.classList.remove('hidden');
          } else {
            item.classList.add('hidden');
          }
        });
      });
    });
  }

  /* ── Animate hero stats counter ───────────────────────── */
  function animateCounter(el, target, suffix) {
    let current = 0;
    const increment = target / 60;
    const timer = setInterval(() => {
      current += increment;
      if (current >= target) {
        clearInterval(timer);
        current = target;
      }
      el.textContent = Math.floor(current) + suffix;
    }, 25);
  }

  const statsObserver = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const counters = entry.target.querySelectorAll('[data-count]');
        counters.forEach(counter => {
          const target = parseInt(counter.dataset.count, 10);
          const suffix = counter.dataset.suffix || '';
          animateCounter(counter, target, suffix);
        });
        statsObserver.unobserve(entry.target);
      }
    });
  }, { threshold: 0.5 });

  const heroStats = document.querySelector('.hero-stats');
  if (heroStats) statsObserver.observe(heroStats);

})();
