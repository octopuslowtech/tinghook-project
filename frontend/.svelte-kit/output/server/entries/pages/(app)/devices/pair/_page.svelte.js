import "clsx";
import "../../../../../chunks/auth.js";
import { U as sanitize_props, V as spread_props, W as slot } from "../../../../../chunks/index2.js";
import { I as Icon } from "../../../../../chunks/Icon.js";
import { S as Smartphone } from "../../../../../chunks/smartphone.js";
function Arrow_left($$renderer, $$props) {
  const $$sanitized_props = sanitize_props($$props);
  /**
   * @license lucide-svelte v0.460.1 - ISC
   *
   * This source code is licensed under the ISC license.
   * See the LICENSE file in the root directory of this source tree.
   */
  const iconNode = [
    ["path", { "d": "m12 19-7-7 7-7" }],
    ["path", { "d": "M19 12H5" }]
  ];
  Icon($$renderer, spread_props([
    { name: "arrow-left" },
    $$sanitized_props,
    {
      /**
       * @component @name ArrowLeft
       * @description Lucide SVG icon component, renders SVG Element with children.
       *
       * @preview ![img](data:image/svg+xml;base64,PHN2ZyAgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIgogIHdpZHRoPSIyNCIKICBoZWlnaHQ9IjI0IgogIHZpZXdCb3g9IjAgMCAyNCAyNCIKICBmaWxsPSJub25lIgogIHN0cm9rZT0iIzAwMCIgc3R5bGU9ImJhY2tncm91bmQtY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDJweCIKICBzdHJva2Utd2lkdGg9IjIiCiAgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIgogIHN0cm9rZS1saW5lam9pbj0icm91bmQiCj4KICA8cGF0aCBkPSJtMTIgMTktNy03IDctNyIgLz4KICA8cGF0aCBkPSJNMTkgMTJINSIgLz4KPC9zdmc+Cg==) - https://lucide.dev/icons/arrow-left
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
    $$renderer2.push(`<a href="/devices" class="mb-6 inline-flex items-center text-gray-600 hover:text-gray-900">`);
    Arrow_left($$renderer2, { class: "mr-2 h-4 w-4" });
    $$renderer2.push(`<!----> Back to Devices</a> <div class="mx-auto max-w-lg"><div class="rounded-lg bg-white p-8 text-center shadow">`);
    Smartphone($$renderer2, { class: "mx-auto mb-4 h-16 w-16 text-blue-500" });
    $$renderer2.push(`<!----> <h1 class="mb-2 text-2xl font-bold">Pair Your Device</h1> <p class="mb-8 text-gray-600">Open the TingHook app on your Android device and scan this QR code.</p> `);
    {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div class="mx-auto h-64 w-64 animate-pulse rounded-lg bg-gray-100"></div>`);
    }
    $$renderer2.push(`<!--]--> <div class="mt-8 rounded-lg bg-gray-50 p-4 text-left"><h3 class="mb-2 font-semibold">Instructions:</h3> <ol class="list-inside list-decimal space-y-2 text-sm text-gray-600"><li>Download TingHook from Google Play Store</li> <li>Open the app and tap "Connect"</li> <li>Scan the QR code above</li> <li>Grant SMS and Notification permissions</li></ol></div></div></div>`);
  });
}
export {
  _page as default
};
