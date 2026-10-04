import { untrack } from "svelte";

/** The current time, updated at each whole minute while a component is subscribed: the page shows
 *  minutes ("prije 3 min", the hour strip, the sky), never seconds. */
function createClock() {
  let now = $state(new Date());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  function tick() {
    if (timer !== null) clearTimeout(timer);
    now = new Date();
    timer = setTimeout(tick, 60_000 - (now.getTime() % 60_000));
  }

  // A phone suspends timers in the background: catch up as soon as the page is shown again.
  function onVisibility() {
    if (document.visibilityState === "visible") tick();
  }

  return {
    get now() {
      return now;
    },
    /** `$effect(() => clock.subscribe())` ties the timer to the component's lifetime. */
    subscribe(): () => void {
      return untrack(() => {
        if (consumers++ === 0) {
          document.addEventListener("visibilitychange", onVisibility);
          tick();
        }
        return () => {
          if (--consumers > 0) return;
          document.removeEventListener("visibilitychange", onVisibility);
          if (timer !== null) clearTimeout(timer);
          timer = null;
        };
      });
    },
  };
}

export const clock = createClock();
