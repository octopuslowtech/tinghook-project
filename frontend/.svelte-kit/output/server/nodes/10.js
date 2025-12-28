

export const index = 10;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/(app)/settings/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/10.DBXJPDJL.js","_app/immutable/chunks/BIHYj5q8.js","_app/immutable/chunks/CNB1M18W.js","_app/immutable/chunks/BO5QB-9v.js"];
export const stylesheets = [];
export const fonts = [];
