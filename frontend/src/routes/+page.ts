import type { PageLoad } from './$types';
import type { WeatherData } from '$lib/types/weather';
import { fetchForecast } from '$svc/forecast';

interface PageData {
  forecast: WeatherData;
  streamed?: {};
}

export const load = (async ({ fetch, params }): Promise<PageData> => {
  const forecast: WeatherData = await fetchForecast(fetch);
  
  return {
    forecast: forecast,
    streamed: {},
  };
}) satisfies PageLoad;