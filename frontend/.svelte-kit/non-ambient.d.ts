
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/(auth)" | "/(app)" | "/" | "/(app)/dashboard" | "/(app)/devices" | "/(app)/devices/pair" | "/(auth)/login" | "/(app)/logs" | "/(auth)/register" | "/(app)/rules" | "/(app)/settings";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/(auth)": Record<string, never>;
			"/(app)": Record<string, never>;
			"/": Record<string, never>;
			"/(app)/dashboard": Record<string, never>;
			"/(app)/devices": Record<string, never>;
			"/(app)/devices/pair": Record<string, never>;
			"/(auth)/login": Record<string, never>;
			"/(app)/logs": Record<string, never>;
			"/(auth)/register": Record<string, never>;
			"/(app)/rules": Record<string, never>;
			"/(app)/settings": Record<string, never>
		};
		Pathname(): "/" | "/dashboard" | "/dashboard/" | "/devices" | "/devices/" | "/devices/pair" | "/devices/pair/" | "/login" | "/login/" | "/logs" | "/logs/" | "/register" | "/register/" | "/rules" | "/rules/" | "/settings" | "/settings/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/favicon.png" | string & {};
	}
}