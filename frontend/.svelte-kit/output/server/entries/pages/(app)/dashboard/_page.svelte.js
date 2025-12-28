import { U as sanitize_props, V as spread_props, W as slot, Y as ensure_array_like } from "../../../../chunks/index2.js";
import { I as Icon } from "../../../../chunks/Icon.js";
function Activity($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.460.1 - ISC
   *
   * This source code is licensed under the ISC license.
   * See the LICENSE file in the root directory of this source tree.
   */
  const iconNode = [
    [
      "path",
      {
        "d": "M22 12h-2.48a2 2 0 0 0-1.93 1.46l-2.35 8.36a.25.25 0 0 1-.48 0L9.24 2.18a.25.25 0 0 0-.48 0l-2.35 8.36A2 2 0 0 1 4.49 12H2"
      }
    ]
  ];
  Icon($$renderer, spread_props([
    { name: "activity" },
    $$sanitized_props,
    {
      /**
       * @component @name Activity
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8cGF0aCBkPSJNMjIgMTJoLTIuNDhhMiAyIDAgMCAwLTEuOTMgMS40NmwtMi4zNSA4LjM2YS4yNS4yNSAwIDAgMS0uNDggMEw5LjI0IDIuMThhLjI1LjI1IDAgMCAwLS40OCAwbC0yLjM1IDguMzZBMiAyIDAgMCAxIDQuNDkgMTJIMiIgLz4KPC9zdmc+Cg==) - https://lucide.dev/icons/activity
       * @see https://lucide.dev/guide/packages/lucide-svelte - Documentation
       *
       * @param {Object} props - Lucide icons props and any valid SVG attribute
       * @returns {FunctionalComponent} Svelte component
       *
       */
      iconNode,
      children: ($$renderer2) => {
        $$renderer2.push(`<!--[-->`);
        slot($$renderer2, $$props, "default", {});
        $$renderer2.push(`<!--]-->`);
      },
      $$slots: { default: true }
    }
  ]));
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    $$renderer2.push(`<h1 class="mb-6 text-2xl font-bold">Dashboard</h1> <div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4"><div class="rounded-lg bg-white p-6 shadow">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="animate-pulse"><div class="mb-2 h-4 w-24 rounded bg-gray-200"></div> <div class="h-8 w-16 rounded bg-gray-200"></div></div>`);
    }
    $$renderer2.push(`<!--]--></div> <div class="rounded-lg bg-white p-6 shadow">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="animate-pulse"><div class="mb-2 h-4 w-24 rounded bg-gray-200"></div> <div class="h-8 w-16 rounded bg-gray-200"></div></div>`);
    }
    $$renderer2.push(`<!--]--></div> <div class="rounded-lg bg-white p-6 shadow">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="animate-pulse"><div class="mb-2 h-4 w-24 rounded bg-gray-200"></div> <div class="h-8 w-16 rounded bg-gray-200"></div></div>`);
    }
    $$renderer2.push(`<!--]--></div> <div class="rounded-lg bg-white p-6 shadow">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="animate-pulse"><div class="mb-2 h-4 w-24 rounded bg-gray-200"></div> <div class="h-8 w-16 rounded bg-gray-200"></div></div>`);
    }
    $$renderer2.push(`<!--]--></div></div> <div class="rounded-lg bg-white shadow"><div class="border-b p-6"><h2 class="flex items-center gap-2 text-lg font-semibold">`);
    Activity($$renderer2, { class: "h-5 w-5" });
    $$renderer2.push(`<!----> Recent Activity</h2></div> <div class="p-6">`);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="animate-pulse space-y-4"><!--[-->`);
      const each_array = ensure_array_like(Array(5));
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        each_array[$$index];
        $$renderer2.push(`<div class="flex items-center gap-4"><div class="h-10 w-10 rounded-full bg-gray-200"></div> <div class="flex-1"><div class="mb-2 h-4 w-48 rounded bg-gray-200"></div> <div class="h-3 w-32 rounded bg-gray-200"></div></div></div>`);
      }
      $$renderer2.push(`<!--]--></div>`);
    }
    $$renderer2.push(`<!--]--></div></div>`);
  });
}
export {
  _page as default
};
