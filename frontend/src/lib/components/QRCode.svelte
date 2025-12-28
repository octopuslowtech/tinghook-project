<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		data: string;
		size?: number;
	}
	let { data, size = 200 }: Props = $props();

	let canvas: HTMLCanvasElement;

	onMount(async () => {
		const QRCodeLib = await import('qrcode');
		QRCodeLib.toCanvas(canvas, data, {
			width: size,
			margin: 2,
			color: { dark: '#000000', light: '#ffffff' }
		});
	});
</script>

<canvas bind:this={canvas}></canvas>
