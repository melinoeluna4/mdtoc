const copyButtons = document.querySelectorAll("[data-copy]");

for (const button of copyButtons) {
  button.addEventListener("click", async () => {
    const value = button.getAttribute("data-copy");
    if (!value) return;

    const original = button.textContent;

    try {
      await navigator.clipboard.writeText(value);
      button.textContent = "Copiado";
    } catch {
      button.textContent = "Copia manualmente";
    }

    window.setTimeout(() => {
      button.textContent = original;
    }, 1600);
  });
}

const tocLinks = [...document.querySelectorAll(".toc a[href^='#']")];
const sections = tocLinks
  .map((link) => document.querySelector(link.getAttribute("href")))
  .filter(Boolean);

if ("IntersectionObserver" in window && sections.length > 0) {
  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        for (const link of tocLinks) {
          const active = link.getAttribute("href") === `#${entry.target.id}`;
          if (active) {
            link.setAttribute("aria-current", "location");
          } else {
            link.removeAttribute("aria-current");
          }
        }
      }
    },
    { rootMargin: "-20% 0px -65% 0px", threshold: 0 },
  );

  for (const section of sections) observer.observe(section);
}
