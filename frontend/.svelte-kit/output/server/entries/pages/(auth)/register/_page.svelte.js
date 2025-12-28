import { a as attr } from "../../../../chunks/attributes.js";
import "@sveltejs/kit/internal";
import "../../../../chunks/exports.js";
import "../../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../../chunks/state.svelte.js";
import "../../../../chunks/auth.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let email = "";
    let password = "";
    let confirmPassword = "";
    let loading = false;
    $$renderer2.push(`<div class="bg-white rounded-lg shadow-md p-8"><h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Create Account</h1> `);
    {
      $$renderer2.push("<!--[!-->");
    }
    $$renderer2.push(`<!--]--> <form><div class="mb-4"><label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label> <input type="email" id="email"${attr("value", email)} class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" placeholder="you@example.com"${attr("disabled", loading, true)}/></div> <div class="mb-4"><label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label> <input type="password" id="password"${attr("value", password)} class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" placeholder="Minimum 8 characters"${attr("disabled", loading, true)}/></div> <div class="mb-6"><label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label> <input type="password" id="confirmPassword"${attr("value", confirmPassword)} class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" placeholder="Confirm your password"${attr("disabled", loading, true)}/></div> <button type="submit"${attr("disabled", loading, true)} class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center">`);
    {
      $$renderer2.push("<!--[!-->");
      $$renderer2.push(`Create Account`);
    }
    $$renderer2.push(`<!--]--></button></form> <p class="mt-4 text-center text-sm text-gray-600">Already have an account? <a href="/login" class="text-blue-600 hover:text-blue-700 font-medium">Sign in</a></p></div>`);
  });
}
export {
  _page as default
};
