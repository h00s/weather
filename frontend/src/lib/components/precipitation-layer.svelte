<script lang="ts">
  import type { Precipitation } from "$lib/helpers/weather";

  type Props = { kind: Exclude<Precipitation, "none"> };
  let { kind }: Props = $props();

  // A fixed pseudo-random field (mulberry32), so drops don't jump around when the layer remounts.
  function random(seed: number) {
    return () => {
      seed = (seed + 0x6d2b79f5) | 0;
      let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
      t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  }

  const count = { drizzle: 28, rain: 56, storm: 72, snow: 48 };
  const snow = $derived(kind === "snow");

  const drops = $derived.by(() => {
    const rnd = random(7);
    return Array.from({ length: count[kind] }, () => ({
      left: rnd() * 104 - 2,
      delay: -rnd() * 12,
      duration: snow ? 9 + rnd() * 7 : 0.7 + rnd() * 0.5,
      size: snow ? 3 + rnd() * 4 : 14 + rnd() * 10,
      opacity: snow ? 0.45 + rnd() * 0.4 : 0.18 + rnd() * 0.22,
    }));
  });
</script>

<div class={["field", { slant: !snow }]} aria-hidden="true">
  {#each drops as d, i (i)}
    <span
      class={snow ? "flake" : "drop"}
      style:left="{d.left}%"
      style:animation-delay="{d.delay}s"
      style:animation-duration="{d.duration}s"
      style:--size="{d.size}px"
      style:opacity={d.opacity}
    ></span>
  {/each}
</div>
{#if kind === "storm"}
  <div class="flash" aria-hidden="true"></div>
{/if}

<style>
  /* Only transform and opacity animate, so the compositor does the work. */
  .field {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }
  .slant {
    transform: rotate(10deg) scale(1.2);
  }
  .drop,
  .flake {
    position: absolute;
    top: -40px;
    will-change: transform;
    animation-timing-function: linear;
    animation-iteration-count: infinite;
  }
  .drop {
    width: 1.5px;
    height: var(--size);
    border-radius: 1px;
    background: linear-gradient(to bottom, transparent, white);
    animation-name: fall;
  }
  .flake {
    width: var(--size);
    height: var(--size);
    border-radius: 50%;
    background: white;
    animation-name: drift;
  }
  @keyframes fall {
    to {
      transform: translateY(110vh);
    }
  }
  @keyframes drift {
    25% {
      transform: translate(12px, 28vh);
    }
    50% {
      transform: translate(-6px, 55vh);
    }
    75% {
      transform: translate(10px, 83vh);
    }
    to {
      transform: translate(0, 110vh);
    }
  }
  /* A faint, rare flash: one pulse every 13 s, far below the photosensitivity limits. */
  .flash {
    position: absolute;
    inset: 0;
    background: white;
    opacity: 0;
    animation: flash 13s ease-out infinite;
  }
  @keyframes flash {
    0%,
    96%,
    100% {
      opacity: 0;
    }
    97% {
      opacity: 0.12;
    }
  }
</style>
