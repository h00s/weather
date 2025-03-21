export function formatTemperature(temp: number): string {
  return `${Math.round(temp)}°C`;
}

export function formatHumidity(humidity: number): string {
  return `${humidity}%`;
}

export function formatPrecipitation(precipitation: number, hideIfZero: boolean = false): string {
  if (hideIfZero && precipitation === 0) return '';
  return `${precipitation}mm`;
}

export function weatherIcon(code: number): string {
  if (code === 0) return '☀️'; // Clear sky
  if (code === 1) return '🌤️'; // Mainly clear
  if (code === 2) return '⛅'; // Partly cloudy
  if (code === 3) return '☁️'; // Overcast
  if (code >= 45 && code <= 48) return '🌫️'; // Fog
  if (code >= 51 && code <= 55) return '🌧️'; // Drizzle
  if (code >= 61 && code <= 65) return '🌧️'; // Rain
  if (code >= 71 && code <= 75) return '❄️'; // Snow
  if (code >= 80 && code <= 82) return '🌦️'; // Showers
  if (code >= 95) return '⛈️'; // Thunderstorm
  return '🌡️'; // Default
}

export function weatherDescription(code: number): string {
  if (code === 0) return 'Clear sky';
  if (code === 1) return 'Mainly clear';
  if (code === 2) return 'Partly cloudy';
  if (code === 3) return 'Overcast';
  if (code >= 45 && code <= 48) return 'Fog';
  if (code >= 51 && code <= 55) return 'Drizzle';
  if (code >= 61 && code <= 65) return 'Rain';
  if (code >= 71 && code <= 75) return 'Snow';
  if (code >= 80 && code <= 82) return 'Rain showers';
  if (code >= 95) return 'Thunderstorm';
  return 'Unknown';
}