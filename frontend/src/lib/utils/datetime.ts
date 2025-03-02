export function formatTimeWithWeekday(timestamp: string): string {
  const date = new Date(timestamp);
  return `${date.toLocaleDateString("en-US", { weekday: "long" })}, ${date.toLocaleTimeString(
    "en-US",
    {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    },
  )}`;
}

export function formatTime(timestamp: string): string {
  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
}